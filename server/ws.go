package server

import (
	"encoding/json"
	"github.com/amanbolat/ca-warehouse-client/logistics"
	bolt "go.etcd.io/bbolt"
	"time"
)

var ShipmentsBucket = []byte("shipments")

func (s *Server) StartShipmentUpdates() {
	go func() {
		ticker := time.Tick(time.Second * 5)
		for {
			select {
			case <-ticker:
				shipments, _, err := s.shipmentStore.GetShipmentUpdates()
				if err != nil {
					s.logger.Errorf("failed to get shipment updates: %v", err)
					continue
				}

				select {
				case s.shipmentsForPrint <- shipments:
				default:
				}

				b, err := json.Marshal(shipments)
				if err != nil {
					s.logger.Errorf("failed to get shipment updates: %v", err)
					continue
				}

				err = s.wsServer.Broadcast(b)
				if err != nil {
					s.logger.Errorf("failed to get shipment updates: %v", err)
				}
			}
		}
	}()
}

func (s *Server) PrintShipmentPrepLabelsOnUpdate() {
	go func() {
		for {
			select {
			case shipments := <-s.shipmentsForPrint:
				var needCheckShipments []string
				var printedShipments []string

				err := s.boltDB.View(func(tx *bolt.Tx) error {
					b := tx.Bucket(ShipmentsBucket)
					for _, sm := range shipments {
						found := b.Get([]byte(sm.Code))
						if found == nil {
							needCheckShipments = append(needCheckShipments, sm.Code)
						}
					}

					return nil
				})

				for _, code := range needCheckShipments {
					s.logger.Debugf("preparation label for %s will be printed", code)
					sm, err := s.shipmentStore.GetShipmentByCode(code)
					if err != nil {
						s.logger.Errorf("failed to print prep labels, %v", err)
						continue
					}

					if sm.CurrentStatusKey == logistics.Preparation {

						label, err := s.labelManger.CreateShipmentPreparationLabels(sm)
						if err != nil {
							s.logger.Errorf("failed to print prep labels, %v", err)
							continue
						}

						err = s.printer.PrintFiles(1, "", label.FullPath)
						if err != nil {
							s.logger.Errorf("failed to print prep labels, %v", err)
							continue
						}
						s.logger.Infof("printing %s shipment preparation label", sm.Code)

						printedShipments = append(printedShipments, sm.Code)
					}
				}

				err = s.boltDB.Update(func(tx *bolt.Tx) error {
					b := tx.Bucket(ShipmentsBucket)
					for _, code := range printedShipments {
						err := b.Put([]byte(code), []byte(code))
						if err != nil {
							return err
						}
					}

					return nil
				})
				if err != nil {
					s.logger.Errorf("failed save printed shipment codes: %v", err)
					continue
				}
			}
		}

	}()
}

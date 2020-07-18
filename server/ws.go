package server

import (
	"encoding/json"
	"time"
)

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

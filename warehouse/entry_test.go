package warehouse_test

import (
	"encoding/json"
	"github.com/amanbolat/ca-warehouse-client/warehouse"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEntryJSON(t *testing.T) {
	str := `{"warehouse":"GZWH2","product_category":"household_goods","has_brand":true,"box_qty":1,"customer_code":"CON","track_code":"123","product_name":"123"}`
	e := warehouse.Entry{}
	err := json.Unmarshal([]byte(str), &e)
	assert.NoError(t, err)
	t.Log(e)
}

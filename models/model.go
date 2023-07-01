package models

import "time"

type Node struct {
	Ip    string             `json:"ip"`
	Enode string             `json:"enode"`
	Specs ComputerSpecs      `json:"specs"`
	Data  []EthereumNodeData `json:"data"`
}

type ComputerSpecs struct {
	CPU          string `json:"cpu"`
	GPU          string `json:"gpu"`
	RAM          int    `json:"ram"`
	Storage      int    `json:"storage"`
	StorageType  string `json:"storage_type"`
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
}

type DBSizeData struct {
	PartName string `json:"part_name"`
	Size     int64  `json:"size"`
}

type BlockProcessingData struct {
	SyncStageName     string  `json:"sync_stage_name"`
	SnapSyncStatus    float64 `json:"snap_sync_status"`
	OldBodiesStatus   float64 `json:"old_bodies_status"`
	OldReceiptsStatus float64 `json:"old_receipts_status"`
	PeersNumber       int64   `json:"peers_number"`
	//Avg               float64 `json:"avg"`
	//StdDev            float64 `json:"std_dev"`
	//Min               float64 `json:"min"`
	//Max               float64 `json:"max"`
}
type HardwareUsageData struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
}

type NetworkUsageData struct {
	NetworkIn  int64 `json:"network_in"`
	NetworkOut int64 `json:"network_out"`
}

// EthereumNodeData
// Db size after sync (with info about each part)
// Block Processing avg, stdDev, min and max on FullySynced node
// Same with BlockProcessing but durting OldBodies and OldReceipts
// Sync time of each sync step for each chain
// Hardware usage during sync stages (to know which consumes what to be able to pinpoint configuration issues on user side)
// Network used during each sync stage
type EthereumNodeData struct {
	Timestamp           time.Time           `json:"timestamp"`
	DBSize              DBSizeData          `json:"db_size"`
	BlockProcessingInfo BlockProcessingData `json:"block_processing_info"`
	HardwareUsageData   HardwareUsageData   `json:"hardware_usage_data"`
	NetworkUsageData    NetworkUsageData    `json:"network_usage_data"`
}

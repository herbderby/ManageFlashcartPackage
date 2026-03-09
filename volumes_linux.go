package main

// listVolumes on Linux returns an empty list. Physical SD card
// volumes are not accessible from the Cowork VM; the tool exists
// so the binary compiles and the MCP server starts. Users working
// in Cowork should use the macOS host for volume operations.
func listVolumes() ([]Volume, error) {
	return nil, nil
}

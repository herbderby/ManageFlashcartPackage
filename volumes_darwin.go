package main

import (
	"os"
	"syscall"
)

// listVolumes enumerates /Volumes/ on macOS and returns metadata
// for each mounted filesystem.
func listVolumes() ([]Volume, error) {
	entries, err := os.ReadDir("/Volumes")
	if err != nil {
		return nil, err
	}

	var volumes []Volume
	for _, entry := range entries {
		mountPoint := "/Volumes/" + entry.Name()
		var stat syscall.Statfs_t
		if err := syscall.Statfs(mountPoint, &stat); err != nil {
			continue
		}

		totalBytes := uint64(stat.Bsize) * stat.Blocks
		freeBytes := uint64(stat.Bsize) * stat.Bavail

		volumes = append(volumes, Volume{
			Name:       entry.Name(),
			MountPoint: mountPoint,
			FSType:     int8SliceToString(stat.Fstypename[:]),
			TotalBytes: totalBytes,
			FreeBytes:  freeBytes,
			TotalHuman: humanBytes(totalBytes),
			FreeHuman:  humanBytes(freeBytes),
		})
	}
	return volumes, nil
}

// int8SliceToString converts a null-terminated int8 slice (as
// returned by Darwin's Statfs_t) to a Go string.
func int8SliceToString(s []int8) string {
	buf := make([]byte, 0, len(s))
	for _, b := range s {
		if b == 0 {
			break
		}
		buf = append(buf, byte(b))
	}
	return string(buf)
}

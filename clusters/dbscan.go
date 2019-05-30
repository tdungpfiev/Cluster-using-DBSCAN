package clusters

import (
	"github.com/willf/bitset"
)

// DBSCAN pseudocode 
// http://en.wikipedia.org/wiki/DBSCAN

// DBSCAN(D, eps, MinPts)
//    C = 0
//    for each unvisited point P in dataset D
//       mark P as visited
//       NeighborPts = regionQuery(P, eps)
//       if sizeof(NeighborPts) < MinPts
//          mark P as NOISE
//       else
//          C = next cluster
//          expandCluster(P, NeighborPts, C, eps, MinPts)

// expandCluster(P, NeighborPts, C, eps, MinPts)
//    add P to cluster C
//    for each point P' in NeighborPts
//       if P' is not visited
//          mark P' as visited
//          NeighborPts' = regionQuery(P', eps)
//          if sizeof(NeighborPts') >= MinPts
//             NeighborPts = NeighborPts joined with NeighborPts'
//       if P' is not yet member of any cluster
//          add P' to cluster C

// regionQuery(P, eps)
//    return all points within P's eps-neighborhood (including P)

type EpsFunction func(pt Point) float64

func DBScan(points PointList, eps float64, minPoints int) (clusters []Cluster, noise []int) {
    // eps is clustering radius in km
	visited := make([]bool, len(points))
	members := make([]bool, len(points))
	clusters = []Cluster{}
	noise = []int{}
	C := 0
	kdTree := NewKDTree(points)

	// SphericalDistanceFast returns distance which is not mutiplied
	// by kRadiusEarth * kDegreeRadian, adjust eps accordingly
	eps = eps / kRadiusEarth / kDegreeRadian

	neighborUnique := bitset.New(uint(len(points)))

	for i := 0; i < len(points); i++ {
		if visited[i] {
			continue
		}
		visited[i] = true

		neighborPts := kdTree.InRange(points[i], eps, nil)
		if len(neighborPts) < minPoints {
			noise = append(noise, i)
		} else {
			cluster := Cluster{C: C, Points: []int{i}}
			members[i] = true
			C++
			// expandCluster goes here inline
			neighborUnique.ClearAll()
			for j := 0; j < len(neighborPts); j++ {
				neighborUnique.Set(uint(neighborPts[j]))
			}

			for j := 0; j < len(neighborPts); j++ {
				k := neighborPts[j]
				if !visited[k] {
					visited[k] = true
					moreNeighbors := kdTree.InRange(points[k], eps, nil)
					if len(moreNeighbors) >= minPoints {
						for _, p := range moreNeighbors {
							if !neighborUnique.Test(uint(p)) {
								neighborPts = append(neighborPts, p)
								neighborUnique.Set(uint(p))
							}
						}
					}
				}

				if !members[k] {
					cluster.Points = append(cluster.Points, k)
					members[k] = true
				}
			}
			clusters = append(clusters, cluster)
		}
	}

	return
}

func RegionQuery(points PointList, P *Point, eps float64) []int {
	result := []int{}

	for i := 0; i < len(points); i++ {
		if points[i].sqDist(P) < eps*eps {
			result = append(result, i)
		}
	}

	return result
}

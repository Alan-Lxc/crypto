package phase

//func GetMsgFromNode(nodeReceive, nodeSend Node, pointmsg point.Pointmsg) {
//	index := pointmsg.GetIndex()
//	log.Printf("Phase 1 :[Node %d] receive point message from [Node %d]\n", nodeReceive.GetLabel, index)
//	p := pointmsg.GetPoint()
//	//Receive the point and store
//	nodeReceive.RecPoint[nodeReceive.RecCounter] = p
//	nodeReceive.RecCounter += 1
//
//	if nodeReceive.RecCounter == nodeReceive.Total {
//		nodeReceive.RecCounter = 0
//		Phase1(nodeReceive)
//	}
//}
//func Phase1(node Node) {
//	log.Printf("[Node %d] now start phase1", node.Label)
//	x_point := make([]*gmp.Int, node.Degree+1)
//	y_point := make([]*gmp.Int, node.Degree+1)
//	for i := 0; i <= node.Degree; i++ {
//		x_point[i] = node.RecPoint[i].X
//		y_point[i] = node.RecPoint[i].Y
//	}
//	p, err := interpolation.LagrangeInterpolate(node.Degree, x_point, y_point, node.P)
//	if err != nil {
//		for i := 0; i <= node.Degree; i++ {
//			log.Print(x_point[i])
//			log.Print(y_point[i])
//		}
//		log.Print(err)
//		panic("Interpolation failed")
//	}
//	node.RecPoly = p
//	log.Printf("Interpolation finished")
//	//Phase2(node)
//}

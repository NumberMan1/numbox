package dag

import (
	"fmt"
	"testing"
)

func TestDAG(t *testing.T) {
	// 创建一个新的 DAG
	dag := NewDAG()

	// 添加节点
	dag.AddNode("A")
	dag.AddNode("B")
	dag.AddNode("C")
	dag.AddNode("D")
	dag.AddNode("E")

	// 添加边
	err := dag.AddConditionalEdge("A", "B", nil)
	if err != nil {
		fmt.Println(err)
	}
	err = dag.AddConditionalEdge("A", "C", []string{"B"})
	if err != nil {
		fmt.Println(err)
	}
	err = dag.AddConditionalEdge("B", "C", []string{"A"})
	if err != nil {
		fmt.Println(err)
	}
	err = dag.AddConditionalEdge("C", "D", nil)
	if err != nil {
		fmt.Println(err)
	}
	err = dag.AddConditionalEdge("D", "E", nil)
	if err != nil {
		fmt.Println(err)
	}

	// 尝试添加一条从 D 到 A 的边，这会导致环的产生，因此应返回错误
	err = dag.AddConditionalEdge("D", "A", nil)
	if err != nil {
		fmt.Println("添加边时出错:", err)
	}

	// 打印 DAG 的结构
	dag.Print()

	// 获取从节点数组 [A, B] 出发直接可到达的所有节点，不包括起始节点本身
	reachable, err := dag.GetDirectlyReachableNodes([]string{"A", "B"})
	if err != nil {
		fmt.Println("获取可到达节点时出错:", err)
	} else {
		fmt.Println("从节点 [A, B] 出发直接可到达的节点:", reachable)
	}

	// 获取从节点数组 [A] 出发直接可到达的所有节点，不包括起始节点本身
	reachable, err = dag.GetDirectlyReachableNodes([]string{"A"})
	if err != nil {
		fmt.Println("获取可到达节点时出错:", err)
	} else {
		fmt.Println("从节点 [A] 出发直接可到达的节点:", reachable)
	}

	// 获取从节点数组 [] 出发直接可到达的所有节点，不包括起始节点本身
	reachable, err = dag.GetDirectlyReachableNodes([]string{})
	if err != nil {
		fmt.Println("获取可到达节点时出错:", err)
	} else {
		fmt.Println("从节点 [] 出发直接可到达的节点:", reachable)
	}
}

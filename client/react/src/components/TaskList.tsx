import { useState, useEffect } from "react";
import axios from "axios";
import { Table, Button } from "react-bootstrap";
import { FaPencil, FaTrash } from "react-icons/fa6";

interface Task {
  id: number;
  name: string;
  due_date: number;
  status: string;
}

function TaskList() {
  const [tasks, setTasks] = useState<Task[]>([]);
  useEffect(() => {
    axios
      .get("http://localhost:3000/api/tasks")
      .then((response) => {
        console.log(response.data);
        let tasks: Task[] = response.data;
        setTasks(tasks);
      })
      .catch((error) => {
        console.error("Error fetching tasks:", error);
      });
  }, []);
  const deleteHandler = (id: number) => {
    axios
      .get("http://localhost:3000/api/delete/task/" + id.toString())
      .then(() => {
        setTasks(tasks.filter((task) => task.id !== id));
      })
      .catch((error) => {
        console.error("Error deleting task:", error);
      });
  };
  return (
    <div>
      <h2>Task List</h2>
      <Table striped bordered hover>
        <thead>
          <tr>
            <th style={{ textAlign: "center" }}>Task Name</th>
            <th style={{ textAlign: "center" }}>Due Date</th>
            <th style={{ textAlign: "center" }}>Status</th>
            <th style={{ width: "110px", textAlign: "center" }}>Actions</th>
          </tr>
        </thead>
        <tbody>
          {tasks.length > 0 ? (
            tasks.map((task) => {
              return (
                <tr key={task.id}>
                  <td>{task.name}</td>
                  <td>{new Date(task.due_date).toLocaleString()}</td>
                  <td>{task.status}</td>
                  <td>
                    <div>
                      <Button variant="primary" style={{ marginRight: "10px" }}>
                        <FaPencil style={{ color: "white" }}></FaPencil>
                      </Button>
                      <Button
                        variant="danger"
                        onClick={() => deleteHandler(task.id)}
                      >
                        <FaTrash style={{ color: "white" }}></FaTrash>
                      </Button>
                    </div>
                  </td>
                </tr>
              );
            })
          ) : (
            <tr>
              <td colSpan={4}>No tasks found</td>
            </tr>
          )}
        </tbody>
      </Table>
    </div>
  );
}

export default TaskList;

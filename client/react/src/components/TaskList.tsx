import { useState, useEffect } from "react";
import axios from "axios";
import { Table, Button, Badge, Stack } from "react-bootstrap";
import { FaTrash } from "react-icons/fa6";
import TaskModals from "./TaskModal";
import AddTaskModal from "./AddTask";
interface Task {
  id: number;
  name: string;
  due_date: number;
  status: number;
}

function TaskList() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const statusArray = ["Unchange", "Pending", "In Progress", "Completed"];
  const variantArray = ["secondary", "primary", "warning", "success"];
  const loadTasks = () => {
    axios
      .get("http://task.localhost/tasks")
      .then((response) => {
        console.log(response.data);
        let tasks: Task[] = response.data;
        setTasks(tasks);
      })
      .catch((error) => {
        console.error("Error fetching tasks:", error);
      });
  };
  useEffect(() => {
    loadTasks();
  }, []);
  const deleteHandler = (id: number) => {
    axios
      .get("http://task.localhost/delete/task/" + id.toString())
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
            <th style={{ textAlign: "center" }}>Due</th>
            <th style={{ textAlign: "center" }}>Status</th>
            <th style={{ width: "90px", textAlign: "center" }}>Actions</th>
          </tr>
        </thead>
        <tbody>
          {tasks.length > 0 ? (
            tasks.map((task) => {
              return (
                <tr key={task.id}>
                  <td>{task.name}</td>
                  <td>{new Date(task.due_date).toLocaleString()}</td>
                  <td>
                    <Badge
                      bg={
                        variantArray[task.status]
                          ? variantArray[task.status]
                          : "secondary"
                      }
                    >
                      {statusArray[task.status]
                        ? statusArray[task.status]
                        : "Bruh"}
                    </Badge>
                  </td>
                  <td>
                    <div>
                      <Stack direction="horizontal" gap={1}>
                        <TaskModals
                          task_id={task.id}
                          reload_parent={loadTasks}
                        ></TaskModals>
                        <Button
                          size="sm"
                          variant="danger"
                          onClick={() => deleteHandler(task.id)}
                        >
                          <FaTrash style={{ color: "white" }}></FaTrash>
                        </Button>
                      </Stack>
                    </div>
                  </td>
                </tr>
              );
            })
          ) : (
            <tr>
              <td colSpan={4}>
                <div style={{ textAlign: "center" }}>No tasks found</div>
              </td>
            </tr>
          )}
        </tbody>
      </Table>
      <AddTaskModal reload_parent={loadTasks}></AddTaskModal>
    </div>
  );
}

export default TaskList;

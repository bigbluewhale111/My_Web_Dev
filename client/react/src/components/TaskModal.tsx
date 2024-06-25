import axios from "axios";
import { useState } from "react";
import {
  Button,
  Modal,
  DropdownButton,
  Dropdown,
  ButtonGroup,
  Badge,
  Form,
  Alert,
  Stack,
} from "react-bootstrap";
import { FaInfo } from "react-icons/fa6";

interface Task {
  id: number;
  name: string;
  description: string;
  due_date: number;
  status: number;
}

function TaskModals(props: { task_id: number; reload_parent: () => void }) {
  const [task, setTask] = useState<Task>();
  const [validDueDate, setValidDueDate] = useState(true);
  const [newName, setNewName] = useState("");
  const [newDescription, setNewDescription] = useState("");
  const [newDueDate, setNewDueDate] = useState("");
  const [newStatus, setNewStatus] = useState(0);
  const statusArray = ["Unchange", "Pending", "In Progress", "Completed"];
  const variantArray = ["secondary", "primary", "warning", "success"];
  const [showInfo, setShowInfo] = useState(false);
  const handleShowInfo = () => {
    axios
      .get("/api/task/" + props.task_id)
      .then((response) => {
        console.log(response.data);
        let task: Task = response.data;
        setTask(task);
      })
      .catch((error) => {
        console.error("Error fetching task:", error);
      });
    setShowInfo(true);
  };
  const handleCloseInfo = () => setShowInfo(false);
  const [showEdit, setShowEdit] = useState(false);
  const handleShowEdit = () => {
    setValidDueDate(true);
    setShowEdit(true);
  };
  const handleCloseEdit = () => setShowEdit(false);
  const handleSendEdit = () => {
    var sendDueDate = Date.parse(newDueDate);
    if (newDueDate != "" && isNaN(sendDueDate)) {
      setValidDueDate(false);
      return;
    }
    if (newDueDate == "") {
      sendDueDate = 0;
    }
    axios
      .post("/api/edit/task/" + props.task_id, {
        name: newName,
        description: newDescription,
        due_date: sendDueDate,
        status: newStatus,
      })
      .then(() => {
        console.log("Task edited successfully");
        props.reload_parent();
        handleCloseEdit();
      })
      .catch((error) => {
        console.error("Error editing task:", error);
      });
  };
  return (
    <>
      <Button size="sm" variant="primary" onClick={handleShowInfo}>
        <FaInfo style={{ color: "white" }}></FaInfo>
      </Button>

      <Modal show={showInfo} onHide={handleCloseInfo} size="lg">
        <Modal.Header closeButton>
          <Modal.Title>
            {task != undefined ? task.name : "No Task Found"}
          </Modal.Title>
        </Modal.Header>
        {task != undefined ? (
          <Modal.Body>
            <Stack gap={3}>
              <Stack direction="horizontal" gap={4}>
                <h4>Due:</h4>
                <h4>{new Date(task.due_date).toLocaleString()}</h4>

                <h4>Status:</h4>
                <h4>
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
                </h4>
              </Stack>

              <p>{task.description}</p>
            </Stack>
          </Modal.Body>
        ) : (
          <></>
        )}
        <Modal.Footer>
          {task != undefined ? (
            <Button
              variant="secondary"
              onClick={() => {
                handleCloseInfo();
                handleShowEdit();
              }}
            >
              Edit
            </Button>
          ) : (
            <></>
          )}

          <Button variant="primary" onClick={handleCloseInfo}>
            Cancel
          </Button>
        </Modal.Footer>
      </Modal>

      <Modal show={showEdit} onHide={handleCloseEdit} size="lg">
        <Modal.Header closeButton>
          <Modal.Title>
            {task != undefined ? "Edit Task" : "No Task Found"}
          </Modal.Title>
        </Modal.Header>
        {task != undefined ? (
          <Modal.Body>
            <Stack gap={3}>
              <Stack direction="horizontal" gap={3}>
                <p>Name:</p>
                <Form.Control
                  placeholder={task.name}
                  onChange={(e) => {
                    setNewName(e.target.value);
                  }}
                ></Form.Control>
              </Stack>
              <Stack direction="horizontal" gap={3}>
                <p>Due:</p>
                <Form.Control
                  placeholder={new Date(task.due_date).toLocaleString()}
                  onChange={(e) => {
                    setNewDueDate(e.target.value);
                    setValidDueDate(true);
                  }}
                ></Form.Control>
                <p>Status:</p>
                <>
                  <DropdownButton
                    as={ButtonGroup}
                    title={statusArray[newStatus]}
                    variant={
                      variantArray[newStatus]
                        ? variantArray[newStatus]
                        : "secondary"
                    }
                  >
                    <Dropdown.Item
                      onClick={() => {
                        setNewStatus(1);
                      }}
                    >
                      Pending
                    </Dropdown.Item>
                    <Dropdown.Item
                      onClick={() => {
                        setNewStatus(2);
                      }}
                    >
                      In Progress
                    </Dropdown.Item>
                    <Dropdown.Item
                      onClick={() => {
                        setNewStatus(3);
                      }}
                    >
                      Completed
                    </Dropdown.Item>
                    <Dropdown.Divider />
                    <Dropdown.Item
                      onClick={() => {
                        setNewStatus(0);
                      }}
                    >
                      Unchange
                    </Dropdown.Item>
                  </DropdownButton>
                </>
              </Stack>
              <Form.Control
                as="textarea"
                placeholder={task.description}
                style={{ height: "100px" }}
                onChange={(e) => {
                  setNewDescription(e.target.value);
                }}
              ></Form.Control>
              {validDueDate ? (
                ""
              ) : (
                <Alert variant="danger">
                  Please enter a valid date in the format MM/DD/YYYY,HH/MM/SS
                  AM/PM
                </Alert>
              )}
            </Stack>
          </Modal.Body>
        ) : (
          <></>
        )}
        <Modal.Footer>
          {task != undefined ? (
            <Button variant="secondary" onClick={handleSendEdit}>
              Edit
            </Button>
          ) : (
            <></>
          )}
          <Button variant="primary" onClick={handleCloseEdit}>
            Close
          </Button>
        </Modal.Footer>
      </Modal>
    </>
  );
}

export default TaskModals;

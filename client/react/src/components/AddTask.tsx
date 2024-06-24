import axios from "axios";
import { useState } from "react";
import {
  Modal,
  Button,
  Form,
  Stack,
  Alert,
  DropdownButton,
  Dropdown,
  ButtonGroup,
} from "react-bootstrap";
import { FaPlus } from "react-icons/fa6";

function AddTaskModal(props: { reload_parent: () => void }) {
  const [showAdd, setShowAdd] = useState(false);
  const handleCloseAdd = () => setShowAdd(false);
  const [validDueDate, setValidDueDate] = useState(true);
  const handleShowAdd = () => setShowAdd(true);
  const [newName, setNewName] = useState("");
  const [newDescription, setNewDescription] = useState("");
  const [newDueDate, setNewDueDate] = useState("");
  const [newStatus, setNewStatus] = useState(1);
  const statusArray = ["Unchange", "Pending", "In Progress", "Completed"];
  const variantArray = ["secondary", "primary", "warning", "success"];
  const handleSendTask = () => {
    var sendDueDate = Date.parse(newDueDate);
    if (newDueDate != "" && isNaN(sendDueDate)) {
      setValidDueDate(false);
      return;
    }
    if (newDueDate == "") {
      sendDueDate = 0;
    }
    axios
      .post("http://task.localhost/create/task", {
        name: newName,
        description: newDescription,
        due_date: sendDueDate,
        status: newStatus,
      })
      .then(() => {
        console.log("Task added successfully");
        props.reload_parent();
        handleCloseAdd();
      })
      .catch((error) => {
        console.error("Error adding task:", error);
      });
  };
  return (
    <>
      <div style={{ textAlign: "center" }}>
        <Button
          variant="primary"
          style={{ justifySelf: "center", alignSelf: "center" }}
          onClick={handleShowAdd}
        >
          <FaPlus style={{ color: "white" }}></FaPlus> New Task
        </Button>
      </div>

      <Modal show={showAdd} onHide={handleCloseAdd} size="lg">
        <Modal.Header closeButton>
          <Modal.Title>Add New Task</Modal.Title>
        </Modal.Header>

        <Modal.Body>
          <Stack gap={3}>
            <Stack direction="horizontal" gap={3}>
              <p>Name:</p>
              <Form.Control
                placeholder="My Task"
                onChange={(e) => {
                  setNewName(e.target.value);
                }}
              ></Form.Control>
            </Stack>
            <Stack direction="horizontal" gap={3}>
              <p>Due:</p>
              <Form.Control
                placeholder={new Date(Date.now()).toLocaleString()}
                className="sm"
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
                </DropdownButton>
              </>
            </Stack>
            <Form.Group controlId="formDescription">
              <Form.Control
                as="textarea"
                placeholder="I dont know what is this task about."
                style={{ height: "100px" }}
                onChange={(e) => {
                  setNewDescription(e.target.value);
                }}
              ></Form.Control>
            </Form.Group>
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
        <Modal.Footer>
          <Button variant="secondary" onClick={handleSendTask}>
            Add
          </Button>
          <Button variant="primary" onClick={handleCloseAdd}>
            Close
          </Button>
        </Modal.Footer>
      </Modal>
    </>
  );
}

export default AddTaskModal;

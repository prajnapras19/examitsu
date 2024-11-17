import React, { useState } from "react";
import { Button, Container, Offcanvas } from "react-bootstrap";

const QuestionListSidebar = (props) => {
  const { questionIDList, handleChooseQuestion } = props;
  const [show, setShow] = useState(false);

  const handleClose = () => setShow(false);
  const handleShow = () => setShow(true);

  return (
    <Container className="prevent-select">
      <Button variant="primary" onClick={handleShow} className="me-2">
        Klik di sini untuk melihat daftar soal
      </Button>
      <Offcanvas show={show} onHide={handleClose} {...props}>
        <Offcanvas.Header closeButton>
          <Offcanvas.Title>Daftar Soal</Offcanvas.Title>
        </Offcanvas.Header>
        <Offcanvas.Body>
          {questionIDList.map((data, i) => (
            <>
              <hr/>
              <Container onClick={() => {
                handleClose();
                handleChooseQuestion(i + 1);
              }}>
                Soal {i + 1}
              </Container>
            </>
          ))}
          <hr/>
        </Offcanvas.Body>
      </Offcanvas>
    </Container>
  );
};

export default QuestionListSidebar;
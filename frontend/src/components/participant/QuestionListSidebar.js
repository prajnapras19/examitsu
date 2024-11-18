import React, { useState } from "react";
import { Button, Container, Offcanvas } from "react-bootstrap";

const QuestionListSidebar = (props) => {
  const { questionIDList, handleChooseQuestion, handleShowSubmitModal } = props;
  const [show, setShow] = useState(false);

  const handleClose = () => setShow(false);
  const handleShow = () => setShow(true);

  return (
    <>
      <Button variant="primary" onClick={handleShow} className="me-3">
        Klik di sini untuk melihat daftar soal
      </Button>
      <Offcanvas show={show} onHide={handleClose} {...props}>
        <Offcanvas.Header closeButton>
          <Offcanvas.Title className="prevent-select">Daftar Soal</Offcanvas.Title>
        </Offcanvas.Header>
        <Offcanvas.Body>
          {questionIDList.map((data, i) => (
            <>
              <hr/>
              <Container className="prevent-select" onClick={() => {
                handleClose();
                handleChooseQuestion(i + 1);
              }}>
                Soal {i + 1}
              </Container>
            </>
          ))}
          <hr/>
          <Container className="text-center mt-5">
            <Button variant="danger" onClick={handleShowSubmitModal}>Kumpulkan ujian</Button>
          </Container>
        </Offcanvas.Body>
      </Offcanvas>
    </>
  );
};

export default QuestionListSidebar;
import React from 'react';
import { Card } from 'react-bootstrap';
import { FaPlus } from "react-icons/fa";

const AddQuestionCard = () => {
  const handleAddQuestion = (examSerial) => {
    // TODO: make empty question
  }

  return (
    <Card className="card" onClick={handleAddQuestion}>
      <Card.Header style={{height: '50%'}}>
        <FaPlus style={{height: '100%'}} size={50}></FaPlus>
      </Card.Header>
      <Card.Body>
        <Card.Title>
          Tambah Soal
        </Card.Title>
        <Card.Text>
          Klik di sini untuk menambah soal ke dalam sistem.
        </Card.Text>
      </Card.Body>
    </Card>
  );
}

export default AddQuestionCard;
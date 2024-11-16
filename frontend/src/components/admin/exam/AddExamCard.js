import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from 'react-bootstrap';
import { FaPlus } from "react-icons/fa";

const AddExamCard = () => {
  const navigate = useNavigate();
  return (
    <Card className="card" onClick={() => navigate('/admin/exams/new')}>
      <Card.Header style={{height: '50%'}}>
        <FaPlus style={{height: '100%'}} size={50}></FaPlus>
      </Card.Header>
      <Card.Body>
        <Card.Title>
          Tambah Ujian
        </Card.Title>
        <Card.Text>
          Klik di sini untuk menambah ujian ke dalam sistem.
        </Card.Text>
      </Card.Body>
    </Card>
  );
}

export default AddExamCard;
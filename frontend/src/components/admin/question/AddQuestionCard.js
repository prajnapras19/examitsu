import React from 'react';
import { Card } from 'react-bootstrap';
import { FaPlus } from "react-icons/fa";
import { useParams } from 'react-router-dom';
import axios from 'axios';
import { toast } from 'react-toastify';

const AddQuestionCard = (props) => {
  const { auth, initiateTriggerRender } = props;
  const { examSerial } = useParams();
  const handleAddQuestion = async () => {
    try {
      await axios.put(`${process.env.REACT_APP_BACKEND_URL}/api/v1/admin/questions`, {
        exam_serial: examSerial,
        data: '{}', // editorjs requirement
      }, {
        headers: {
          Authorization: `Bearer ${auth.token}`,
        },
      });
  
      toast.success('Soal berhasil ditambahkan!', {
        position: "top-center",
        autoClose: 3000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
        draggable: true,
      });
    } catch (err) {
      toast.error('Gagal menambahkan soal. Silakan coba beberapa saat lagi.', {
        position: "top-center",
        autoClose: 5000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
        draggable: true,
      });
    }
    
    initiateTriggerRender();
  }

  return (
    <Card className="card" onClick={() => handleAddQuestion()}>
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
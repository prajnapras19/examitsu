import React, { useState } from 'react';
import { Card } from 'react-bootstrap';
import { FaUpload } from "react-icons/fa";
import { toast } from 'react-toastify';
import axios from 'axios';

const UploadExamFileCard = ({auth, onClick, onFinish}) => {
  const handleFileChange = (event) => {
    onClick();
    handleUpload(event.target.files[0]);
  };

  const handleUpload = (file) => {
    const formData = new FormData();
    formData.append("file", file);
    axios.post(`${process.env.REACT_APP_BACKEND_URL}/api/v1/admin/exams/upload`, formData, {
      headers: {
        "Content-Type": "multipart/form-data",
        'Authorization': `Bearer ${auth.token}`
      },
    })
    .then(response => {
      toast.success('Ujian berhasil diunggah!', {
        position: "top-center",
        autoClose: 3000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
        draggable: true,
      });
      onFinish();
    })
    .catch(err => {
      if (err.code === 500) {
        toast.error(`Sedang terjadi masalah pada server. Silakan coba beberapa saat lagi.`, {
          position: "top-center",
          autoClose: 5000,
          hideProgressBar: false,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
        });
      } else {
        toast.error(`Gagal memproses pengunggahan. Silakan coba menggunakan data yang lain.`, {
          position: "top-center",
          autoClose: 5000,
          hideProgressBar: false,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
        });
      }
      onFinish();
    });
  };

  const triggerFileInput = () => {
    document.getElementById("fileInput").click();
  };
  return (
    <Card className="card" onClick={triggerFileInput}>
      <input
        type="file"
        id="fileInput"
        style={{ display: "none" }}
        onChange={handleFileChange}
      />
      <Card.Header style={{height: '50%'}}>
        <FaUpload style={{height: '100%'}} size={50}></FaUpload>
      </Card.Header>
      <Card.Body>
        <Card.Title>
          Unggah Ujian
        </Card.Title>
        <Card.Text>
          Klik di sini untuk mengunggah ujian.
        </Card.Text>
      </Card.Body>
    </Card>
  );
}

export default UploadExamFileCard;
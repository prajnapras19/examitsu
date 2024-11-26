import React from 'react';
import { Card } from 'react-bootstrap';
import { FaDownload } from "react-icons/fa";
import { toast } from 'react-toastify';
import axios from 'axios';

const DownloadExamTemplateCard = ({auth}) => {
  const handleDownload = async () => {
    try {
      const response = await axios.get(`${process.env.REACT_APP_BACKEND_URL}/api/v1/admin/exams/template`, {
        responseType: "blob",
        headers: {
          'Authorization': `Bearer ${auth.token}`
        },
      });

      const blob = new Blob([response.data], { type: "application/zip" });

      const contentDisposition = response.headers["Content-Disposition"];
      const fileName = contentDisposition
        ? contentDisposition
            .split("filename=")[1]
            .replace(/["']/g, "") // Remove any quotes around the filename
        : "example.zip";

      const link = document.createElement("a");
      link.download = fileName;
      link.href = window.URL.createObjectURL(blob);
      link.click();

      window.URL.revokeObjectURL(link.href);
    } catch (error) {
      toast.error(`Sedang terjadi masalah pada server. Silakan coba beberapa saat lagi.`, {
        position: "top-center",
        autoClose: 5000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
        draggable: true,
      });
    }
  };
  return (
    <Card className="card" onClick={handleDownload}>
      <Card.Header style={{height: '50%'}}>
        <FaDownload style={{height: '100%'}} size={50}></FaDownload>
      </Card.Header>
      <Card.Body>
        <Card.Title>
          Unduh Format untuk Mengunggah Soal
        </Card.Title>
        <Card.Text>
          Klik di sini untuk mengunduh format soal untuk diunggah.
        </Card.Text>
      </Card.Body>
    </Card>
  );
}

export default DownloadExamTemplateCard;
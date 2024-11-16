import React, { useEffect, useState, useRef } from 'react';
import axios from 'axios';
import { Modal, Button } from "react-bootstrap";
import EditorJS from "@editorjs/editorjs";
import { toast } from 'react-toastify';

const EditQuestionModal = (props) => {
    const {show, onClose, questionId, auth} = props;
    const [loading, setLoading] = useState(false);
    const editorInstance = useRef(null);

    const initializeEditor = async () => {
      try {
        const response = await axios.post(`${process.env.REACT_APP_BACKEND_URL}/api/v1/admin/questions/${questionId}`,
          {}, {
            headers: {
              Authorization: `Bearer ${auth.token}`,
            },
          },
        );
        // TODO: image
        editorInstance.current = new EditorJS({
          holder: "editor",
          tools: {
            header: {
              class: require("@editorjs/header"),
              inlineToolbar: ["link", "bold", "italic"],
            },
            list: require("@editorjs/list"),
            underline: require("@editorjs/underline"),
          },
          data: JSON.parse(response.data.data.data),
        });
      } catch (err) {
        console.log('err', err);
        toast.error(`Gagal mendapatkan data soal. Silakan coba beberapa saat lagi.`, {
          position: "top-center",
          autoClose: 5000,
          hideProgressBar: false,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
        });
        handleClose();
      }
    };

    const handleSubmit = async () => {
      if (!editorInstance.current) return;
  
      setLoading(true);
      try {
        const outputData = await editorInstance.current.save();
        await axios.patch(`${process.env.REACT_APP_BACKEND_URL}/api/v1/admin/questions/${questionId}`, {
            data: JSON.stringify(outputData),
          }, {
            headers: {
              'Authorization': `Bearer ${auth.token}`
            },
          },
        );
        toast.success('Soal berhasil diubah!', {
          position: "top-center",
          autoClose: 3000,
          hideProgressBar: false,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
        });
        onClose(); // Close the modal
      } catch (error) {
        toast.error(`Gagal mengubah soal. Silakan coba beberapa saat lagi.`, {
          position: "top-center",
          autoClose: 5000,
          hideProgressBar: false,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
        });
      } finally {
        setLoading(false);
      }
    };

    const handleClose = () => {
      editorInstance.current?.destroy();
      onClose();
    };

    useEffect(() => {
      if (show) {
        initializeEditor();
      }
    }, [show]);

    return (
      <Modal show={show} onHide={handleClose} size="lg" centered>
        <Modal.Header closeButton>
          <Modal.Title>Ubah Soal</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <h3>Soal</h3>
          <div id="editor" style={{ border: "1px solid #ccc", minHeight: "200px", padding: "10px" }}></div>
          <h3>Pilihan Jawaban</h3>
          <Button variant="primary" disabled={loading}>
            Tambah Pilihan Jawaban
          </Button>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={handleClose} disabled={loading}>
            Batal
          </Button>
          <Button variant="primary" onClick={handleSubmit} disabled={loading}>
            {loading ? "Menyimpan..." : "Simpan"}
          </Button>
        </Modal.Footer>
      </Modal>
    );
}

export default EditQuestionModal;
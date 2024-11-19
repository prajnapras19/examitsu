import React, { useEffect, useState, useRef } from 'react';
import axios from 'axios';
import { Form, Modal, Button } from "react-bootstrap";
import EditorJS from "@editorjs/editorjs";
import Paragraph from "@editorjs/paragraph";
import Underline from "@editorjs/underline";
import List from "@editorjs/list";
import Header from "@editorjs/header";
import { toast } from 'react-toastify';
import ImageTool from '@editorjs/image';

const EditQuestionModal = (props) => {
    const {show, onClose, questionId, auth} = props;
    const [loading, setLoading] = useState(false);
    const editorInstance = useRef(null);
    const [mcqOptions, setMcqOptions] = useState([]);

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
              class: Header,
              inlineToolbar: ["bold", "italic", "underline"],
            },
            list: List,
            underline: {
              class: Underline,
              shortcut: 'CTRL+U',
            },
            paragraph: {
              class: Paragraph,
              inlineToolbar: ["bold", "italic", "underline"],
            },
            image: {
              class: ImageTool,
              config: {
                uploader: {
                  async uploadByFile(file) {
                    // Fetch the signed URL
                    const getFileUploadURL = await fetch(`${process.env.REACT_APP_BACKEND_URL}/api/v1/admin/questions/file-upload-url`, {
                      method: 'POST',
                      headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${auth.token}`,
                      },
                      body: JSON.stringify({ file_type: file.type }),
                    });
        
                    if (!getFileUploadURL.ok) {
                      throw new Error('Gagal mendapatkan alamat pengunggahan!');
                    }
        
                    const res = await getFileUploadURL.json();
                    const { upload_url, public_url } = res.data;
                    console.log('res', res.data);
        
                    // Upload file to the signed URL
                    const uploadResponse = await fetch(upload_url, {
                      method: 'PUT',
                      headers: { 'Content-Type': file.type },
                      body: file,
                    });
        
                    if (!uploadResponse.ok) {
                      throw new Error('Gagal mengunggah gambar!');
                    }
        
                    // Return the public URL for Editor.js to use
                    return {
                      success: 1,
                      file: {
                        url: public_url, // URL accessible by your application or users
                      },
                    };
                  },
                },
              },
            },
          },
          data: JSON.parse(response.data.data.data),
        });
      } catch (err) {
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
        const fetchData = async () => {
          try {
            const response = await axios.post(`${process.env.REACT_APP_BACKEND_URL}/api/v1/admin/mcq-options/question-id/${questionId}`,
              {}, {
                headers: {
                  Authorization: `Bearer ${auth.token}`,
                },
              },
            );
            setMcqOptions(response.data.data);
          } catch (err) {
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
        fetchData();
      }
    }, [show]);

    const handleOnChangeMcqOptions = (e, idx) => {
      const {name, value} = e.target;

      if (name === 'point' && value === '') {
        return;
      }

      const currentObject = mcqOptions[idx];

      if (name === 'point') {
        currentObject[name] = Number(value);
      } else {
        currentObject[name] = value;
      }
      
      
      setMcqOptions(
        (prevArray) =>
          prevArray.map((item, i) => (i === idx ? currentObject : item))
      );
    }

    const handleSaveMcqOption = async (e, idx) => {
      e.preventDefault();
      
      try {
        await axios.patch(`${process.env.REACT_APP_BACKEND_URL}/api/v1/admin/mcq-options/${mcqOptions[idx].id}`, mcqOptions[idx], {
          headers: {
            Authorization: `Bearer ${auth.token}`,
          },
        });
  
        toast.success('Pilihan jawaban berhasil diubah!', {
          position: "top-center",
          autoClose: 3000,
          hideProgressBar: false,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
        });
      } catch (err) {
        toast.error(`Gagal mengubah pilihan jawaban, silakan coba beberapa saat lagi.`, {
          position: "top-center",
          autoClose: 5000,
          hideProgressBar: false,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
        });
      }
    }

    const handleAddMcqOption = async () => {
      try {
        const response = await axios.put(`${process.env.REACT_APP_BACKEND_URL}/api/v1/admin/mcq-options`, {
          question_id: questionId,
        }, {
          headers: {
            Authorization: `Bearer ${auth.token}`,
          },
        });

        setMcqOptions((prev) => prev.concat(response.data.data));
    
        toast.success('Pilihan jawaban berhasil ditambahkan!', {
          position: "top-center",
          autoClose: 3000,
          hideProgressBar: false,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
        });
      } catch (err) {
        toast.error('Gagal menambahkan pilihan jawaban. Silakan coba beberapa saat lagi.', {
          position: "top-center",
          autoClose: 5000,
          hideProgressBar: false,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
        });
      }
    }

    const handleDeleteMcqOption = async (idx) => {
      try {
        await axios.delete(`${process.env.REACT_APP_BACKEND_URL}/api/v1/admin/mcq-options/${mcqOptions[idx].id}`, {
          headers: {
            Authorization: `Bearer ${auth.token}`,
          },
        });

        setMcqOptions(
          (prevArray) => 
            prevArray.filter((_, index) => index !== idx)
        );
    
        toast.success('Pilihan jawaban berhasil dihapus!', {
          position: "top-center",
          autoClose: 3000,
          hideProgressBar: false,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
        });
      } catch (err) {
        toast.error('Gagal menghapus pilihan jawaban. Silakan coba beberapa saat lagi.', {
          position: "top-center",
          autoClose: 5000,
          hideProgressBar: false,
          closeOnClick: true,
          pauseOnHover: true,
          draggable: true,
        });
      }
    }

    return (
      <Modal show={show} onHide={handleClose} size="lg" centered>
        <Modal.Header closeButton>
          <Modal.Title>Ubah Soal</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <h3>Soal</h3>
          <div id="editor" style={{ border: "1px solid #ccc", minHeight: "200px", padding: "10px" }}></div>
          <br/>
          <Button variant="primary" onClick={handleSubmit} disabled={loading}>
            {loading ? "Menyimpan..." : "Simpan"}
          </Button>
          <hr/>
          <h3>Pilihan Jawaban</h3>
          <hr/>
          {mcqOptions.map((mcqOption, i) => (
            <>
              <Form className="my-4" onSubmit={(e) => handleSaveMcqOption(e, i)}>
                <Form.Group className="my-3" controlId="description" key="description">
                  <Form.Label><b>Deskripsi</b></Form.Label>
                  <Form.Control
                    type='text'
                    name='description'
                    value={mcqOption.description}
                    onChange={(e) => handleOnChangeMcqOptions(e, i)}
                    autoComplete='off'
                  />
                </Form.Group>
                <Form.Group className="my-3" controlId="point" key="point">
                  <Form.Label><b>Poin</b></Form.Label>
                  <Form.Control
                    type='number'
                    name='point'
                    step='1'
                    value={mcqOption.point}
                    onChange={(e) => handleOnChangeMcqOptions(e, i)}
                    autoComplete='off'
                  />
                </Form.Group>
                <Button variant="primary" type="submit" className="me-3">Simpan</Button>
                <Button variant="danger" onClick={() => handleDeleteMcqOption(i)}>Hapus</Button>
              </Form>
              <hr/>
            </>
          ))}
          <Button variant="primary" onClick={handleAddMcqOption}>
            Tambah Pilihan Jawaban
          </Button>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={handleClose}>
            Tutup
          </Button>
        </Modal.Footer>
      </Modal>
    );
}

export default EditQuestionModal;
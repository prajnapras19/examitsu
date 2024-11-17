import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate, useParams } from 'react-router-dom';
import { toast } from 'react-toastify';
import { Form, Container, Spinner, Button } from 'react-bootstrap';
import BackToHomepageCard from '../home/BackToHomepageCard';
import ReadExamsMenuCard from '../exam/ReadExamsMenuCard';
import ReadParticipantsOfThisExamMenuCard from './ReadParticipantsOfThisExamMenuCard';

const AddParticipants = (props) => {
  const { auth } = props;
  const navigate = useNavigate();
  const { examSerial } = useParams();
  const [disableSubmitButton, setDisableSubmitButton] = useState(false);

  const fields = [
    {
      label: 'Nama-nama Peserta',
      placeholder: 'Masukkan Nama-nama Peserta',
      name: 'names',
      type: 'textarea',
      rows: 8,
      required: true,
      defaultValue: '',
    },
    {
      label: 'Kata Sandi (jika dibiarkan kosong, maka kata sandi akan dibuat oleh sistem)',
      placeholder: 'Masukkan Kata Sandi',
      name: 'password',
      type: 'text',
      required: false,
      defaultValue: '',
    },
  ]

  const [formData, setFormData] = useState(
    fields.reduce((acc, field) => ({ ...acc, [field.name]: field.defaultValue }), {})
  );

  const handleInputChange = (e) => {
    const { name, type } = e.target;
    if (type === 'checkbox') {
      const { checked } = e.target;
      setFormData({ ...formData, [name]: checked });
    } else {
      const { value } = e.target;
      setFormData({ ...formData, [name]: value });
    }
  };

  useEffect(() => {
    if (auth.loading) {
      return;
    }
    if (!auth.isLoggedIn) {
      navigate('/admin/login');
    }
  }, [auth.loading, auth.isLoggedIn]);

  if (auth.loading) {
    return (
      <Container className="text-center">
        <Spinner animation="border" />
        <p>Mohon tunggu...</p>
      </Container>
    );
  }

  const handleSubmit = async (e) => {
    e.preventDefault();

    const customObject = {
      exam_serial: examSerial,
      names: formData.names.split('\n'),
      password: formData.password,
    };

    console.log('customObject', customObject);

    try {
      await axios.put(`${process.env.REACT_APP_BACKEND_URL}/api/v1/admin/participants`, customObject, {
        headers: {
          Authorization: `Bearer ${auth.token}`,
        },
      });

      toast.success('Peserta berhasil ditambahkan!', {
        position: "top-center",
        autoClose: 3000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
        draggable: true,
      });
      navigate(`/admin/exams/${examSerial}/participants`);
    } catch (err) {
      toast.error(`Gagal menambahkan peserta, silakan coba beberapa saat lagi atau menggunakan data yang lain.`, {
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
    <Container>
      <h1 className="my-4">Tambah Peserta</h1>
      <hr/>
      <Container className="text-center mt-5">
        <Container className="card-grid">
          <BackToHomepageCard></BackToHomepageCard>
          <ReadExamsMenuCard></ReadExamsMenuCard>
          <ReadParticipantsOfThisExamMenuCard></ReadParticipantsOfThisExamMenuCard>
        </Container>
      </Container>
      <hr/>

      <Form className="my-4" onSubmit={handleSubmit}>
        {fields.map((field) => (
          <Form.Group className="my-3" controlId={field.name} key={field.name}>
            <Form.Label><b>{field.label}{field.required ? <span> (harus diisi)</span> : ''}</b></Form.Label>
            {
              field.type === 'float'
              ? (
                <Form.Control
                  type='number'
                  min={field.minValue}
                  step={field.step}
                  name={field.name}
                  value={formData[field.name]}
                  onChange={handleInputChange}
                  required={field.required}
                  placeholder={field.placeholder}
                  autoComplete='off'
                />
              )
              : field.type === 'boolean'
              ? (
                <Form.Check
                  type="switch"
                  name={field.name}
                  checked={formData[field.name]}
                  onChange={handleInputChange}
                />
              )
              : field.type === 'textarea'
              ? (
                <Form.Control
                  as='textarea'
                  rows={field.rows}
                  type={field.type}
                  name={field.name}
                  value={formData[field.name]}
                  onChange={handleInputChange}
                  required={field.required}
                  placeholder={field.placeholder}
                  autoComplete='off'
                />
              )
              : (
                <Form.Control
                  type={field.type}
                  name={field.name}
                  value={formData[field.name]}
                  onChange={handleInputChange}
                  required={field.required}
                  placeholder={field.placeholder}
                  autoComplete='off'
                />
              )
            }
          </Form.Group>
        ))}
        <Button variant="primary" type="submit" className="mt-3">
          Kirim
        </Button>
      </Form>

    </Container>
  );
}

export default AddParticipants;
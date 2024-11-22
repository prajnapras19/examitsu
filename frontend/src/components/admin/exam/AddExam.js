import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';
import { Form, Container, Spinner, Button } from 'react-bootstrap';
import ReadExamsMenuCard from './ReadExamsMenuCard';
import BackToHomepageCard from '../home/BackToHomepageCard';

const AddExam = (props) => {
  const { auth } = props;
  const navigate = useNavigate();

  const fields = [
    {
      label: 'Nama Ujian',
      name: 'name',
      type: 'text',
      required: true,
      defaultValue: '',
    },
    {
      label: 'Sudah Bisa Dikerjakan?',
      name: 'is_open',
      type: 'boolean',
      defaultValue: false,
    },
    {
      label: 'Durasi Pengerjaan (dalam satuan menit, durasi ini dapat diubah untuk setiap peserta jika dibutuhkan)',
      name: 'allowed_duration_minutes',
      type: 'float',
      minValue: 1,
      defaultValue: 120,
      step: 1,
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
      if (name === 'allowed_duration_minutes' && value === '') {
        return;
      }
      if (name === 'allowed_duration_minutes') {
        setFormData({ ...formData, [name]: Number(value) });
      } else {
        setFormData({ ...formData, [name]: value });
      }
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
      ...formData
    };

    try {
      await axios.put(`${process.env.REACT_APP_BACKEND_URL}/api/v1/admin/exams`, customObject, {
        headers: {
          Authorization: `Bearer ${auth.token}`,
        },
      });

      toast.success('Ujian berhasil ditambahkan!', {
        position: "top-center",
        autoClose: 3000,
        hideProgressBar: false,
        closeOnClick: true,
        pauseOnHover: true,
        draggable: true,
      });
      navigate('/admin/exams');
    } catch (err) {
      toast.error(`Gagal menambahkan ujian, silakan coba menggunakan data yang lain.`, {
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
      <h1 className="my-4">Tambah Ujian</h1>
      <hr/>
      <Container className="text-center mt-5">
        <Container className="card-grid">
          <BackToHomepageCard></BackToHomepageCard>
          <ReadExamsMenuCard></ReadExamsMenuCard>
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
                  placeholder={`Masukkan ${field.label}`}
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
              : (
                <Form.Control
                  type={field.type}
                  name={field.name}
                  value={formData[field.name]}
                  onChange={handleInputChange}
                  required={field.required}
                  placeholder={`Masukkan ${field.label}`}
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

export default AddExam;
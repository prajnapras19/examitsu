import { Route, Routes } from 'react-router-dom';
import useAdminAuth from './hooks/useAdminAuth';
import Login from './components/admin/auth/Login';
import Homepage from './components/admin/home/Homepage';
import NotFoundPage from './components/etc/404';
import ReadExams from './components/admin/exam/ReadExams';
import AddExam from './components/admin/exam/AddExam';
import EditExam from './components/admin/exam/EditExam';

function AdminRoutes() {
  const adminAuth = useAdminAuth();
  return (
    <Routes>
      <Route path="/login" element={<Login auth={adminAuth}/>}/>
      <Route path="/home" element={<Homepage auth={adminAuth}/>}/>
      <Route path="/exams" element={<ReadExams auth={adminAuth} />} />
      <Route path="/exams/new" element={<AddExam auth={adminAuth} />} />
      <Route path="/exams/edit/:examSerial" element={<EditExam auth={adminAuth} />} />
      <Route path="*" element={<NotFoundPage/>}/>
    </Routes>
  );
}

export default AdminRoutes;
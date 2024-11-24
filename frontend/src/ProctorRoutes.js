import { Route, Routes } from 'react-router-dom';
import useProctorAuth from './hooks/useProctorAuth';
import NotFoundPage from './components/etc/404';
import AuthorizeSession from './components/proctor/AuthorizeSession';
import Login from './components/proctor/auth/Login';

function ProctorRoutes() {
  const proctorAuth = useProctorAuth();
  return (
    <Routes>
      <Route path="/login" element={<Login auth={proctorAuth}/>}/>
      <Route path="/authorize" element={<AuthorizeSession auth={proctorAuth} />} />
      <Route path="*" element={<NotFoundPage/>}/>
    </Routes>
  );
}

export default ProctorRoutes;
import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'react-toastify/dist/ReactToastify.css';
import { ToastContainer } from 'react-toastify';
import NotFoundPage from './components/etc/404';
import ScrollToTop from './components/etc/ScrollToTop';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import useAdminAuth from './hooks/useAdminAuth';
import Login from './components/admin/auth/Login';
import Homepage from './components/admin/home/Homepage';

function App() {
  const adminAuth = useAdminAuth();
  return (
    <div className="App my-4">
      <Router>
        <ScrollToTop />
        <ToastContainer />
        <Routes>
          <Route path="/admin/login" element={<Login auth={adminAuth}/>}/>
          <Route path="/admin/home" element={<Homepage auth={adminAuth}/>}/>
          <Route path="*" element={<NotFoundPage/>}/>
        </Routes>
      </Router>
    </div>
  );
}

export default App;

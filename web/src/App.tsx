import { BrowserRouter, Routes, Route } from "react-router-dom";
import { AuthProvider } from "./context/AuthContext";
import { ProtectedRoute } from "./components/ProtectedRoute";
import { Layout } from "./components/Layout";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Dashboard from "./pages/Dashboard";
import ChangePassword from "./pages/ChangePassword";
import Permissions from "./pages/Permissions";
import Roles from "./pages/Roles";
import RoleForm from "./pages/RoleForm";
import Locations from "./pages/Locations";
import LocationForm from "./pages/LocationForm";
import Users from "./pages/Users";
import UserForm from "./pages/UserForm";
import Patients from "./pages/Patients";
import PatientForm from "./pages/PatientForm";
import PatientDetail from "./pages/PatientDetail";

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/" element={<ProtectedRoute><Layout /></ProtectedRoute>}>
            <Route index element={<Dashboard />} />
            <Route path="change-password" element={<ChangePassword />} />
            <Route path="permissions" element={<Permissions />} />
            <Route path="roles" element={<Roles />} />
            <Route path="roles/:id" element={<RoleForm />} />
            <Route path="locations" element={<Locations />} />
            <Route path="locations/:id" element={<LocationForm />} />
            <Route path="users" element={<Users />} />
            <Route path="users/:id" element={<UserForm />} />
            <Route path="patients" element={<Patients />} />
            <Route path="patients/new" element={<PatientForm />} />
            <Route path="patients/:id" element={<PatientDetail />} />
            <Route path="patients/:id/edit" element={<PatientForm />} />
          </Route>
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}

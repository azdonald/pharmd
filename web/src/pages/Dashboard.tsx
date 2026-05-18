import { useAuth } from "../context/AuthContext";

export default function Dashboard() {
  const { user } = useAuth();
  return (
    <div>
      <h1>Dashboard</h1>
      <p>Welcome, {user?.first_name} {user?.last_name}</p>
      <p>Organisation: {user?.organisation_name}</p>
    </div>
  );
}

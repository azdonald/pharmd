import { NavLink, Outlet, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export function Layout() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  return (
    <div className="layout">
      <aside className="sidebar">
        <div className="sidebar-header">
          <h2>PharmD</h2>
        </div>
        <nav>
          <NavLink to="/">Dashboard</NavLink>
          <NavLink to="/users">Users</NavLink>
          <NavLink to="/roles">Roles</NavLink>
          <NavLink to="/permissions">Permissions</NavLink>
          <NavLink to="/locations">Locations</NavLink>
          <NavLink to="/patients">Patients</NavLink>
          <NavLink to="/products">Products</NavLink>
          <NavLink to="/categories">Categories</NavLink>
          <NavLink to="/inventory">Inventory</NavLink>
          <NavLink to="/suppliers">Suppliers</NavLink>
          <NavLink to="/purchases">Purchases</NavLink>
          <NavLink to="/pricing">Pricing</NavLink>
          <NavLink to="/discounts">Discounts</NavLink>
        </nav>
        <div className="sidebar-footer">
          <span>{user?.first_name} {user?.last_name}</span>
          <span className="org-name">{user?.organisation_name}</span>
          <button onClick={handleLogout}>Logout</button>
        </div>
      </aside>
      <main className="main-content">
        <Outlet />
      </main>
    </div>
  );
}

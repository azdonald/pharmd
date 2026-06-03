import { NavLink, Outlet, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import {
  LayoutDashboard, Users, ShieldCheck, Key, MapPin, Stethoscope,
  Package, Tags, Warehouse, Truck, ShoppingCart, DollarSign,
  Percent, UserCheck, FileText, Pill, CreditCard, BarChart3,
  LogOut,
} from "lucide-react";

const navItems = [
  { to: "/", label: "Dashboard", icon: LayoutDashboard },
  { to: "/users", label: "Users", icon: Users },
  { to: "/roles", label: "Roles", icon: ShieldCheck },
  { to: "/permissions", label: "Permissions", icon: Key },
  { to: "/locations", label: "Locations", icon: MapPin },
  { to: "/patients", label: "Patients", icon: Stethoscope },
  { to: "/products", label: "Products", icon: Package },
  { to: "/categories", label: "Categories", icon: Tags },
  { to: "/inventory", label: "Inventory", icon: Warehouse },
  { to: "/suppliers", label: "Suppliers", icon: Truck },
  { to: "/purchases", label: "Purchases", icon: ShoppingCart },
  { to: "/pricing", label: "Pricing", icon: DollarSign },
  { to: "/discounts", label: "Discounts", icon: Percent },
  { to: "/prescribers", label: "Prescribers", icon: UserCheck },
  { to: "/prescriptions", label: "Prescriptions", icon: FileText },
  { to: "/dispensing", label: "Dispensing", icon: Pill },
  { to: "/pos", label: "POS", icon: CreditCard },
  { to: "/sales", label: "Sales", icon: BarChart3 },
];

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
          {navItems.map(({ to, label, icon: Icon }) => (
            <NavLink key={to} to={to} end={to === "/"}>
              <Icon size={18} />
              <span>{label}</span>
            </NavLink>
          ))}
        </nav>
        <div className="sidebar-footer">
          <div className="sidebar-user">
            <div className="sidebar-user-avatar">
              {user?.first_name?.[0]}{user?.last_name?.[0]}
            </div>
            <div>
              <span className="sidebar-user-name">{user?.first_name} {user?.last_name}</span>
              <span className="org-name">{user?.organisation_name}</span>
            </div>
          </div>
          <button className="sidebar-logout" onClick={handleLogout}>
            <LogOut size={16} />
            Logout
          </button>
        </div>
      </aside>
      <main className="main-content">
        <Outlet />
      </main>
    </div>
  );
}

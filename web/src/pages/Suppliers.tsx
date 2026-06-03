import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { listSuppliers, deleteSupplier, type Supplier } from "../api/suppliers";
import { useToast } from "../context/ToastContext";

export default function Suppliers() {
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const { showToast } = useToast();

  const load = () => {
    setLoading(true);
    listSuppliers(page, 20, search).then(res => {
      setSuppliers(res.data);
      setTotal(res.total);
    }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, [page, search]);

  const handleDelete = async (id: string) => {
    if (!confirm("Delete this supplier?")) return;
    try {
      await deleteSupplier(id);
      load();
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Delete failed", "error");
    }
  };

  const totalPages = Math.ceil(total / 20);

  return (
    <div>
      <div className="page-header">
        <h1>Suppliers</h1>
        <Link to="/suppliers/new" className="btn">New Supplier</Link>
      </div>

      <form className="search-bar" style={{ marginBottom: 12 }}>
        <input value={search} onChange={e => { setSearch(e.target.value); setPage(1); }} placeholder="Search suppliers..." />
        <button type="submit" onClick={e => e.preventDefault()}>Search</button>
      </form>

      {loading ? <p>Loading...</p> : (
        <>
          <table>
            <thead>
              <tr><th>Name</th><th>Contact</th><th>Phone</th><th>Email</th><th>Active</th><th>Actions</th></tr>
            </thead>
            <tbody>
              {suppliers.map(s => (
                <tr key={s.id}>
                  <td><Link to={`/suppliers/${s.id}`}>{s.name}</Link></td>
                  <td>{s.contact_person || "—"}</td>
                  <td>{s.phone || "—"}</td>
                  <td>{s.email || "—"}</td>
                  <td><span className={`badge ${s.is_active ? "badge-active" : "badge-inactive"}`}>{s.is_active ? "Active" : "Inactive"}</span></td>
                  <td>
                    <Link to={`/suppliers/${s.id}/edit`} className="action-link">Edit</Link>
                    <button onClick={() => handleDelete(s.id)} className="action-link action-link-danger">Delete</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>

          {totalPages > 1 && (
            <div className="pagination">
              {page > 1 && <button onClick={() => setPage(page - 1)}>Previous</button>}
              <span>Page {page} of {totalPages}</span>
              {page < totalPages && <button onClick={() => setPage(page + 1)}>Next</button>}
            </div>
          )}
        </>
      )}
    </div>
  );
}

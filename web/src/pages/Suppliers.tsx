import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { listSuppliers, deleteSupplier, type Supplier } from "../api/suppliers";

export default function Suppliers() {
  const [suppliers, setSuppliers] = useState<Supplier[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);

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
      alert(err instanceof Error ? err.message : "Delete failed");
    }
  };

  const totalPages = Math.ceil(total / 20);

  return (
    <div>
      <div className="page-header">
        <h1>Suppliers</h1>
        <Link to="/suppliers/new" className="btn">New Supplier</Link>
      </div>

      <input value={search} onChange={e => { setSearch(e.target.value); setPage(1); }} placeholder="Search suppliers..." style={{ width: "100%", marginBottom: 16 }} />

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
                  <td>{s.is_active ? "Yes" : "No"}</td>
                  <td>
                    <Link to={`/suppliers/${s.id}/edit`}>Edit</Link>
                    <button onClick={() => handleDelete(s.id)}>Delete</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>

          {totalPages > 1 && (
            <div style={{ marginTop: 16, display: "flex", gap: 8 }}>
              {page > 1 && <button onClick={() => setPage(page - 1)}>Previous</button>}
              <span style={{ padding: "8px 0" }}>Page {page} of {totalPages}</span>
              {page < totalPages && <button onClick={() => setPage(page + 1)}>Next</button>}
            </div>
          )}
        </>
      )}
    </div>
  );
}

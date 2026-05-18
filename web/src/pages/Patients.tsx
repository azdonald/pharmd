import { useEffect, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { listPatients, deletePatient, type Patient } from "../api/patients";

export default function Patients() {
  const [patients, setPatients] = useState<Patient[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [searchParams, setSearchParams] = useSearchParams();
  const [search, setSearch] = useState(searchParams.get("query") || "");

  const page = parseInt(searchParams.get("page") || "1");
  const query = searchParams.get("query") || "";

  const load = () => {
    setLoading(true);
    listPatients(page, 20, query).then(res => {
      setPatients(res.data);
      setTotal(res.total);
    }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, [page, query]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setSearchParams({ query: search, page: "1" });
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Deactivate this patient?")) return;
    try {
      await deletePatient(id);
      load();
    } catch (err) {
      alert(err instanceof Error ? err.message : "Delete failed");
    }
  };

  const totalPages = Math.ceil(total / 20);

  return (
    <div>
      <div className="page-header">
        <h1>Patients</h1>
        <Link to="/patients/new" className="btn">New Patient</Link>
      </div>

      <form onSubmit={handleSearch} style={{ display: "flex", gap: 8, marginBottom: 16, maxWidth: "100%" }}>
        <input
          value={search}
          onChange={e => setSearch(e.target.value)}
          placeholder="Search by name, phone, or email..."
          style={{ flex: 1 }}
        />
        <button type="submit">Search</button>
      </form>

      {loading ? <p>Loading...</p> : (
        <>
          <table>
            <thead>
              <tr><th>Name</th><th>Phone</th><th>Email</th><th>Gender</th><th>Active</th><th>Actions</th></tr>
            </thead>
            <tbody>
              {patients.map(p => (
                <tr key={p.id}>
                  <td><Link to={`/patients/${p.id}`}>{p.first_name} {p.last_name}</Link></td>
                  <td>{p.phone}</td>
                  <td>{p.email}</td>
                  <td>{p.gender}</td>
                  <td>{p.is_active ? "Yes" : "No"}</td>
                  <td>
                    <Link to={`/patients/${p.id}/edit`}>Edit</Link>
                    <button onClick={() => handleDelete(p.id)}>Deactivate</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>

          {totalPages > 1 && (
            <div style={{ marginTop: 16, display: "flex", gap: 8 }}>
              {page > 1 && (
                <button onClick={() => setSearchParams({ query, page: String(page - 1) })}>Previous</button>
              )}
              <span style={{ padding: "8px 0" }}>Page {page} of {totalPages}</span>
              {page < totalPages && (
                <button onClick={() => setSearchParams({ query, page: String(page + 1) })}>Next</button>
              )}
            </div>
          )}
        </>
      )}
    </div>
  );
}

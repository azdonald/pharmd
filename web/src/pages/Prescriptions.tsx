import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { listPrescriptions, type Prescription } from "../api/prescriptions";

export default function Prescriptions() {
  const [rxs, setRxs] = useState<Prescription[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [status, setStatus] = useState("");
  const [page, setPage] = useState(1);

  const load = () => {
    setLoading(true);
    listPrescriptions(page, 20, status).then(res => {
      setRxs(res.data);
      setTotal(res.total);
    }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, [page, status]);

  const totalPages = Math.ceil(total / 20);

  const statusClass = (s: string) => {
    switch (s) {
      case "active": return "badge badge-approved";
      case "filled": return "badge badge-received";
      case "expired": return "badge badge-rejected";
      default: return "badge";
    }
  };

  return (
    <div>
      <div className="page-header">
        <h1>Prescriptions</h1>
        <Link to="/prescriptions/new" className="btn">New Prescription</Link>
      </div>

      <select value={status} onChange={e => { setStatus(e.target.value); setPage(1); }} style={{ marginBottom: 16 }}>
        <option value="">All statuses</option>
        <option value="active">Active</option>
        <option value="filled">Filled</option>
        <option value="expired">Expired</option>
      </select>

      {loading ? <p>Loading...</p> : (
        <>
          <table>
            <thead>
              <tr><th>Patient</th><th>Prescriber</th><th>Status</th><th>Issued</th><th>Actions</th></tr>
            </thead>
            <tbody>
              {rxs.map(rx => (
                <tr key={rx.id}>
                  <td>{rx.patient_name}</td>
                  <td>{rx.prescriber_name}</td>
                  <td><span className={statusClass(rx.status)}>{rx.status}</span></td>
                  <td>{rx.issued_date || "—"}</td>
                  <td><Link to={`/prescriptions/${rx.id}`}>View</Link></td>
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

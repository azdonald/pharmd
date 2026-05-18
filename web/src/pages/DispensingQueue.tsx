import { useEffect, useState } from "react";
import { listDispensing, type DispenseRecord } from "../api/dispensing";

export default function DispensingQueue() {
  const [records, setRecords] = useState<DispenseRecord[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [status, setStatus] = useState("");
  const [page, setPage] = useState(1);

  const load = () => {
    setLoading(true);
    listDispensing(page, 20, status).then(res => {
      setRecords(res.data);
      setTotal(res.total);
    }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, [page, status]);

  const totalPages = Math.ceil(total / 20);

  return (
    <div>
      <div className="page-header">
        <h1>Dispensing Queue</h1>
      </div>

      <select value={status} onChange={e => { setStatus(e.target.value); setPage(1); }} style={{ marginBottom: 16 }}>
        <option value="">All statuses</option>
        <option value="pending">Pending</option>
        <option value="in_progress">In Progress</option>
        <option value="completed">Completed</option>
        <option value="cancelled">Cancelled</option>
      </select>

      {loading ? <p>Loading...</p> : (
        <>
          <table>
            <thead>
              <tr><th>Patient</th><th>Product</th><th>Pharmacist</th><th>Status</th><th>Dispensed</th><th>Controlled</th><th>Actions</th></tr>
            </thead>
            <tbody>
              {records.map(dr => (
                <tr key={dr.id}>
                  <td>{dr.patient_name}</td>
                  <td>{dr.product_name}</td>
                  <td>{dr.pharmacist_name}</td>
                  <td><span className={`badge badge-${dr.status}`}>{dr.status}</span></td>
                  <td>{dr.dispensed_at ? new Date(dr.dispensed_at).toLocaleDateString() : "—"}</td>
                  <td>{dr.is_controlled ? "Yes" : "No"}</td>
                  <td><a href={`#/dispensing/${dr.id}`}>View</a></td>
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

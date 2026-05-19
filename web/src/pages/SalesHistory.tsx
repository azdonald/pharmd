import { useEffect, useState } from "react";
import { listSales, updateSale, getReceipt, type Sale } from "../api/pos";

export default function SalesHistory() {
  const [sales, setSales] = useState<Sale[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [status, setStatus] = useState("");
  const [page, setPage] = useState(1);

  const load = () => {
    setLoading(true);
    listSales(page, 20, status).then(res => {
      setSales(res.data);
      setTotal(res.total);
    }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, [page, status]);

  const totalPages = Math.ceil(total / 20);

  const handleVoid = async (id: string) => {
    if (!confirm("Void this sale?")) return;
    try { await updateSale(id, { status: "voided" }); load(); }
    catch (err) { alert(err instanceof Error ? err.message : "Void failed"); }
  };

  const handleRefund = async (id: string) => {
    if (!confirm("Refund this sale?")) return;
    try { await updateSale(id, { status: "refunded" }); load(); }
    catch (err) { alert(err instanceof Error ? err.message : "Refund failed"); }
  };

  const handleReceipt = async (id: string) => {
    try {
      const receipt = await getReceipt(id);
      alert(JSON.stringify(receipt, null, 2));
    } catch (err) {
      alert(err instanceof Error ? err.message : "Receipt failed");
    }
  };

  return (
    <div>
      <div className="page-header">
        <h1>Sales History</h1>
      </div>

      <select value={status} onChange={e => { setStatus(e.target.value); setPage(1); }} style={{ marginBottom: 16 }}>
        <option value="">All statuses</option>
        <option value="active">Active</option>
        <option value="completed">Completed</option>
        <option value="held">Held</option>
        <option value="voided">Voided</option>
        <option value="refunded">Refunded</option>
      </select>

      {loading ? <p>Loading...</p> : (
        <>
          <table>
            <thead>
              <tr><th>Sale #</th><th>Type</th><th>Patient</th><th>Status</th><th>Total</th><th>Paid</th><th>Date</th><th>Actions</th></tr>
            </thead>
            <tbody>
              {sales.map(s => (
                <tr key={s.id}>
                  <td>{s.sale_type === "prescription" ? "Rx" : "OTC"}</td>
                  <td>{s.id.slice(0, 8)}</td>
                  <td>{s.patient_name || "—"}</td>
                  <td><span className={`badge badge-${s.status}`}>{s.status}</span></td>
                  <td>${Number(s.grand_total).toFixed(2)}</td>
                  <td>${Number(s.paid_amount).toFixed(2)}</td>
                  <td>{s.created_at ? new Date(s.created_at).toLocaleDateString() : "—"}</td>
                  <td>
                    <button onClick={() => handleReceipt(s.id)}>Receipt</button>
                    {s.status === "completed" && <button onClick={() => handleRefund(s.id)}>Refund</button>}
                    {(s.status === "active" || s.status === "completed") && <button onClick={() => handleVoid(s.id)}>Void</button>}
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

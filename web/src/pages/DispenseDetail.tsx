import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { getDispenseRecord, createDispense, updateDispenseStatus, checkInteractions, checkAllergies, getLabelData, type DispenseRecord } from "../api/dispensing";
import { getPrescription, type Prescription } from "../api/prescriptions";
import { listUsers, type User } from "../api/users";
import { useToast } from "../context/ToastContext";

export default function DispenseDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [dr, setDr] = useState<DispenseRecord | null>(null);
  const [loading, setLoading] = useState(true);
  const [rx, setRx] = useState<Prescription | null>(null);
  const [users, setUsers] = useState<User[]>([]);
  const [interactions, setInteractions] = useState<any[]>([]);
  const [allergies, setAllergies] = useState<any[]>([]);
  const [showDispense, setShowDispense] = useState(false);
  const { showToast } = useToast();
  const [formItemId, setFormItemId] = useState("");
  const [formQty, setFormQty] = useState("1");
  const [formPharmacist, setFormPharmacist] = useState("");
  const [formTechnician, setFormTechnician] = useState("");
  const [formWitness, setFormWitness] = useState("");
  const [formControlled, setFormControlled] = useState(false);
  const [formNotes, setFormNotes] = useState("");

  const load = () => {
    if (!id) return;
    setLoading(true);
    getDispenseRecord(id)
      .then(drRes => {
        setDr(drRes);
        const rxPromise = drRes.prescription_id ? getPrescription(drRes.prescription_id) : Promise.resolve(null);
        const intPromise = drRes.product_id && drRes.patient_id ? checkInteractions(drRes.product_id, drRes.patient_id) : Promise.resolve([]);
        const allPromise = drRes.product_id && drRes.patient_id ? checkAllergies(drRes.product_id, drRes.patient_id) : Promise.resolve([]);
        return Promise.all([rxPromise, listUsers(1, 100), intPromise, allPromise]);
      })
      .then(([rxRes, usersRes, intRes, allRes]) => {
        setRx(rxRes as any);
        setUsers((usersRes as any).data);
        setInteractions(intRes as any[]);
        setAllergies(allRes as any[]);
      }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, [id]);

  const handleStatusChange = async (status: string) => {
    if (!id) return;
    try {
      const updated = await updateDispenseStatus(id, status);
      setDr(updated);
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Status update failed", "error");
    }
  };

  const handleDispense = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formItemId || !formQty || !formPharmacist) return;
    try {
      const created = await createDispense({
        prescription_item_id: formItemId,
        quantity_dispensed: Number(formQty),
        pharmacist_id: formPharmacist,
        technician_id: formTechnician || undefined,
        witness_name: formWitness || undefined,
        notes: formNotes || undefined,
        is_controlled: formControlled,
      });
      navigate(`/app/dispensing/${created.id}`);
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Dispense failed", "error");
    }
  };

  if (loading) return <p>Loading...</p>;
  if (!dr) return <p>Dispense record not found</p>;

  return (
    <div>
      <div className="page-header">
        <h1>Dispense Record</h1>
        <button onClick={() => navigate("/app/dispensing")}>Back</button>
      </div>

      <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 16, marginBottom: 24 }}>
        <div>
          <p><strong>Patient:</strong> {dr.patient_name}</p>
          <p><strong>Product:</strong> {dr.product_name}</p>
          <p><strong>Pharmacist:</strong> {dr.pharmacist_name}</p>
          <p><strong>Status:</strong> <span className={`badge badge-${dr.status}`}>{dr.status}</span></p>
        </div>
        <div>
          <p><strong>Qty Dispensed:</strong> {dr.quantity_dispensed} / {dr.quantity_prescribed}</p>
          <p><strong>Dispensed At:</strong> {dr.dispensed_at || "—"}</p>
          <p><strong>Controlled:</strong> {dr.is_controlled ? "Yes" : "No"}</p>
          <p><strong>Witness:</strong> {dr.witness_name || "—"}</p>
          <p><strong>Notes:</strong> {dr.notes || "—"}</p>
        </div>
      </div>

      {interactions.length > 0 && (
        <div style={{ marginBottom: 16, padding: 12, border: "1px solid #dc3545", borderRadius: 4, background: "#fff5f5" }}>
          <h3>Drug Interaction Warnings</h3>
          {interactions.map((w, i) => <p key={i}><strong>{w.severity}:</strong> {w.message}</p>)}
        </div>
      )}

      {allergies.length > 0 && (
        <div style={{ marginBottom: 16, padding: 12, border: "1px solid #ffc107", borderRadius: 4, background: "#fffcf0" }}>
          <h3>Allergy Warnings</h3>
          {allergies.map((w, i) => <p key={i}><strong>{w.severity}:</strong> {w.allergen} — {w.reaction || "reaction"}</p>)}
        </div>
      )}

      <div style={{ marginTop: 16, display: "flex", gap: 8 }}>
        {dr.status === "pending" && <button onClick={() => handleStatusChange("in_progress")}>Start Dispensing</button>}
        {dr.status === "in_progress" && <button onClick={() => handleStatusChange("completed")}>Complete</button>}
        {dr.status !== "cancelled" && <button onClick={() => handleStatusChange("cancelled")}>Cancel</button>}
      </div>

      {!showDispense && (
        <button onClick={() => setShowDispense(true)} style={{ marginTop: 16 }}>New Dispense (from Rx)</button>
      )}

      {showDispense && (
        <form onSubmit={handleDispense} style={{ marginTop: 16, padding: 16, border: "1px solid #ddd", borderRadius: 4 }}>
          <h3>Dispense Medication</h3>
          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 8, marginBottom: 8 }}>
            <div>
              <label>Prescription Item</label>
              <select value={formItemId} onChange={e => setFormItemId(e.target.value)} required>
                <option value="">Select item</option>
                {(rx?.items || []).map(item => (
                  <option key={item.id} value={item.id}>{item.product_name} — {item.dosage} {item.frequency}</option>
                ))}
              </select>
            </div>
            <div>
              <label>Quantity *</label>
              <input type="number" min={1} value={formQty} onChange={e => setFormQty(e.target.value)} required />
            </div>
            <div>
              <label>Pharmacist *</label>
              <select value={formPharmacist} onChange={e => setFormPharmacist(e.target.value)} required>
                <option value="">Select pharmacist</option>
                {users.map(u => <option key={u.id} value={u.id}>{u.first_name} {u.last_name}</option>)}
              </select>
            </div>
            <div>
              <label>Technician</label>
              <select value={formTechnician} onChange={e => setFormTechnician(e.target.value)}>
                <option value="">None</option>
                {users.map(u => <option key={u.id} value={u.id}>{u.first_name} {u.last_name}</option>)}
              </select>
            </div>
            <div>
              <label>Witness Name (controlled substances)</label>
              <input value={formWitness} onChange={e => setFormWitness(e.target.value)} placeholder="Required for controlled" />
            </div>
            <div>
              <label>
                <input type="checkbox" checked={formControlled} onChange={e => setFormControlled(e.target.checked)} />
                {" "}Controlled substance
              </label>
            </div>
            <div style={{ gridColumn: "1 / -1" }}>
              <label>Notes</label>
              <textarea value={formNotes} onChange={e => setFormNotes(e.target.value)} rows={2} />
            </div>
          </div>
          <button type="submit">Dispense</button>
          <button type="button" onClick={() => setShowDispense(false)}>Cancel</button>
        </form>
      )}

      {dr.id && (
        <div style={{ marginTop: 16 }}>
          <button onClick={async () => {
            const data = await getLabelData(dr.id);
            showToast(JSON.stringify(data, null, 2).slice(0, 200));
          }}>Print Label</button>
        </div>
      )}
    </div>
  );
}

export function DispenseForm() {
  return <DispenseDetail />;
}

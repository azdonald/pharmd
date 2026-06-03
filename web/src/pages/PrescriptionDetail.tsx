import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { getPrescription, createPrescription, updatePrescription, recordRefill, deletePrescription, type Prescription } from "../api/prescriptions";
import { listPatients, type Patient } from "../api/patients";
import { listLocations, type Location } from "../api/locations";
import { listProducts, type Product } from "../api/products";
import { listPrescribers, type Prescriber } from "../api/prescribers";

export default function PrescriptionDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [rx, setRx] = useState<Prescription | null>(null);
  const [patients, setPatients] = useState<Patient[]>([]);
  const [prescribers, setPrescribers] = useState<Prescriber[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);

  const load = () => {
    if (!id) return;
    setLoading(true);
    getPrescription(id)
      .then(rxRes => Promise.all([rxRes, listPatients(1, 200), listPrescribers(1, 200), listProducts(1, 200)]))
      .then(([rxRes, patRes, prRes, prodRes]) => {
        setRx(rxRes);
        setPatients(patRes.data);
        setPrescribers(prRes.data);
        setProducts(prodRes.data);
      }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, [id]);

  const handleRefill = async (itemId: string) => {
    if (!id) return;
    try {
      const updated = await recordRefill(id, itemId);
      setRx(updated);
    } catch (err) {
      alert(err instanceof Error ? err.message : "Refill failed");
    }
  };

  const handleStatusChange = async (status: string) => {
    if (!id) return;
    try {
      const updated = await updatePrescription(id, { status });
      setRx(updated);
    } catch (err) {
      alert(err instanceof Error ? err.message : "Update failed");
    }
  };

  const handleDelete = async () => {
    if (!id || !confirm("Delete this prescription?")) return;
    try {
      await deletePrescription(id);
      navigate("/prescriptions");
    } catch (err) {
      alert(err instanceof Error ? err.message : "Delete failed");
    }
  };

  const getPatientName = (id: string) => patients.find(p => p.id === id)?.first_name || id;
  const getPrescriberName = (id: string) => prescribers.find(p => p.id === id)?.name || id;
  const getProductName = (id: string) => products.find(p => p.id === id)?.name || id;

  if (loading) return <p>Loading...</p>;
  if (!rx) return <p>Prescription not found</p>;

  return (
    <div>
      <div className="page-header">
        <h1>Prescription</h1>
        <button onClick={() => navigate("/prescriptions")}>Back</button>
      </div>

      <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 16, marginBottom: 24 }}>
        <div>
          <p><strong>Patient:</strong> {rx.patient_name || getPatientName(rx.patient_id)}</p>
          <p><strong>Prescriber:</strong> {rx.prescriber_name || getPrescriberName(rx.prescriber_id)}</p>
          <p><strong>Status:</strong> <span className={`badge badge-${rx.status}`}>{rx.status}</span></p>
        </div>
        <div>
          <p><strong>Diagnosis:</strong> {rx.diagnosis || "—"}</p>
          <p><strong>Issued:</strong> {rx.issued_date || "—"}</p>
          <p><strong>Expires:</strong> {rx.expiry_date || "—"}</p>
          <p><strong>Notes:</strong> {rx.notes || "—"}</p>
        </div>
      </div>

      <h3>Items</h3>
      <table>
        <thead>
          <tr><th>Product</th><th>Dosage</th><th>Frequency</th><th>Qty</th><th>Refills</th><th>Actions</th></tr>
        </thead>
        <tbody>
          {(rx.items || []).map(item => (
            <tr key={item.id}>
              <td>{getProductName(item.product_id)}</td>
              <td>{item.dosage}</td>
              <td>{item.frequency}</td>
              <td>{item.quantity}</td>
              <td>{item.refills_used}/{item.refills_authorized}</td>
              <td>
                {item.refills_used < item.refills_authorized && (
                  <button onClick={() => handleRefill(item.id)}>Refill</button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      <div style={{ marginTop: 16, display: "flex", gap: 8 }}>
        {rx.status === "active" && <button onClick={() => handleStatusChange("filled")}>Mark Filled</button>}
        {rx.status !== "expired" && <button onClick={() => handleStatusChange("expired")}>Mark Expired</button>}
        <button onClick={handleDelete}>Delete</button>
      </div>
    </div>
  );
}

export function PrescriptionForm() {
  const navigate = useNavigate();
  const [patients, setPatients] = useState<Patient[]>([]);
  const [prescribers, setPrescribers] = useState<Prescriber[]>([]);
  const [locations, setLocations] = useState<Location[]>([]);
  const [products, setProducts] = useState<Product[]>([]);

  const [patientId, setPatientId] = useState("");
  const [prescriberId, setPrescriberId] = useState("");
  const [locationId, setLocationId] = useState("");
  const [diagnosis, setDiagnosis] = useState("");
  const [notes, setNotes] = useState("");
  const [issuedDate, setIssuedDate] = useState("");
  const [expiryDate, setExpiryDate] = useState("");
  const [items, setItems] = useState<{ product_id: string; dosage: string; frequency: string; duration: string; quantity: number; refills_authorized: number; notes: string }[]>([]);

  useEffect(() => {
    Promise.all([
      listPatients(1, 200),
      listPrescribers(1, 200),
      listLocations(1, 100),
      listProducts(1, 200),
    ]).then(([patRes, prRes, locRes, prodRes]) => {
      setPatients(patRes.data);
      setPrescribers(prRes.data);
      setLocations(locRes.data);
      setProducts(prodRes.data);
    }).catch(console.error);
  }, []);

  const addItem = () => {
    setItems([...items, { product_id: "", dosage: "", frequency: "", duration: "", quantity: 1, refills_authorized: 0, notes: "" }]);
  };

  const updateItem = (idx: number, field: string, value: any) => {
    const updated = [...items];
    (updated[idx] as any)[field] = value;
    setItems(updated);
  };

  const removeItem = (idx: number) => {
    setItems(items.filter((_, i) => i !== idx));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!patientId || !prescriberId || !locationId || items.length === 0) return;
    if (items.some(i => !i.product_id || !i.dosage || !i.frequency)) { alert("Each item needs product, dosage and frequency"); return; }
    try {
      const created = await createPrescription({
        patient_id: patientId,
        prescriber_id: prescriberId,
        location_id: locationId,
        diagnosis: diagnosis || undefined,
        notes: notes || undefined,
        issued_date: issuedDate || undefined,
        expiry_date: expiryDate || undefined,
        items: items.map(i => ({ ...i, duration: i.duration || undefined, notes: i.notes || undefined })),
      });
      navigate(`/prescriptions/${created.id}`);
    } catch (err) {
      alert(err instanceof Error ? err.message : "Create failed");
    }
  };

  return (
    <div>
      <div className="page-header">
        <h1>New Prescription</h1>
      </div>

      <div className="form-page">
      <form onSubmit={handleSubmit}>
        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 16, marginBottom: 16 }}>
          <div>
            <label>Patient *</label>
            <select value={patientId} onChange={e => setPatientId(e.target.value)} required>
              <option value="">Select patient</option>
              {patients.map(p => <option key={p.id} value={p.id}>{p.first_name} {p.last_name}</option>)}
            </select>
          </div>
          <div>
            <label>Prescriber *</label>
            <select value={prescriberId} onChange={e => setPrescriberId(e.target.value)} required>
              <option value="">Select prescriber</option>
              {prescribers.map(p => <option key={p.id} value={p.id}>{p.name}</option>)}
            </select>
          </div>
          <div>
            <label>Location *</label>
            <select value={locationId} onChange={e => setLocationId(e.target.value)} required>
              <option value="">Select location</option>
              {locations.map(l => <option key={l.id} value={l.id}>{l.name}</option>)}
            </select>
          </div>
          <div>
            <label>Issued Date</label>
            <input type="date" value={issuedDate} onChange={e => setIssuedDate(e.target.value)} />
          </div>
          <div>
            <label>Expiry Date</label>
            <input type="date" value={expiryDate} onChange={e => setExpiryDate(e.target.value)} />
          </div>
          <div>
            <label>Diagnosis</label>
            <textarea value={diagnosis} onChange={e => setDiagnosis(e.target.value)} rows={2} />
          </div>
          <div>
            <label>Notes</label>
            <textarea value={notes} onChange={e => setNotes(e.target.value)} rows={2} />
          </div>
        </div>

        <h3>Items</h3>
        <table>
          <thead>
            <tr><th>Product</th><th>Dosage</th><th>Frequency</th><th>Duration</th><th>Qty</th><th>Refills</th><th></th></tr>
          </thead>
          <tbody>
            {items.map((item, idx) => (
              <tr key={idx}>
                <td><select value={item.product_id} onChange={e => updateItem(idx, "product_id", e.target.value)} required>
                  <option value="">Select product</option>
                  {products.map(p => <option key={p.id} value={p.id}>{p.name}</option>)}
                </select></td>
                <td><input value={item.dosage} onChange={e => updateItem(idx, "dosage", e.target.value)} required placeholder="e.g. 500mg" /></td>
                <td><input value={item.frequency} onChange={e => updateItem(idx, "frequency", e.target.value)} required placeholder="e.g. BID" /></td>
                <td><input value={item.duration} onChange={e => updateItem(idx, "duration", e.target.value)} placeholder="e.g. 7 days" /></td>
                <td><input type="number" min={1} value={item.quantity} onChange={e => updateItem(idx, "quantity", Number(e.target.value))} /></td>
                <td><input type="number" min={0} value={item.refills_authorized} onChange={e => updateItem(idx, "refills_authorized", Number(e.target.value))} /></td>
                <td><button type="button" onClick={() => removeItem(idx)}>Remove</button></td>
              </tr>
            ))}
          </tbody>
        </table>
        <button type="button" onClick={addItem} style={{ marginTop: 8 }}>Add Item</button>

        <div style={{ marginTop: 16, display: "flex", gap: 8 }}>
          <button type="submit">Create Prescription</button>
          <button type="button" onClick={() => navigate("/prescriptions")}>Cancel</button>
        </div>
      </form>
      </div>
    </div>
  );
}

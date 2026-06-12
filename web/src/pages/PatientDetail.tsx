import { useEffect, useState } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import {
  getPatient, deletePatient,
  listPatientAllergies, addPatientAllergy, removePatientAllergy,
  listPatientConditions, addPatientCondition,
  type Patient, type PatientAllergy, type PatientCondition,
} from "../api/patients";
import { useToast } from "../context/ToastContext";

export default function PatientDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [patient, setPatient] = useState<Patient | null>(null);
  const [allergies, setAllergies] = useState<PatientAllergy[]>([]);
  const [conditions, setConditions] = useState<PatientCondition[]>([]);
  const [newAllergy, setNewAllergy] = useState("");
  const [newSeverity, setNewSeverity] = useState("");
  const [newCondition, setNewCondition] = useState("");
  const [loading, setLoading] = useState(true);
  const { showToast } = useToast();

  const load = () => {
    if (!id) return;
    setLoading(true);
    Promise.all([
      getPatient(id),
      listPatientAllergies(id),
      listPatientConditions(id),
    ]).then(([p, a, c]) => {
      setPatient(p);
      setAllergies(a);
      setConditions(c);
    }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, [id]);

  const handleDelete = async () => {
    if (!confirm("Deactivate this patient?")) return;
    try {
      await deletePatient(id!);
      navigate("/app/patients");
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Delete failed", "error");
    }
  };

  const handleAddAllergy = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newAllergy.trim()) return;
    try {
      await addPatientAllergy(id!, { allergy: newAllergy, severity: newSeverity || undefined });
      setNewAllergy("");
      setNewSeverity("");
      const a = await listPatientAllergies(id!);
      setAllergies(a);
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Failed", "error");
    }
  };

  const handleRemoveAllergy = async (allergyId: string) => {
    try {
      await removePatientAllergy(id!, allergyId);
      setAllergies(allergies.filter(a => a.id !== allergyId));
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Failed", "error");
    }
  };

  const handleAddCondition = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newCondition.trim()) return;
    try {
      await addPatientCondition(id!, { condition: newCondition });
      setNewCondition("");
      const c = await listPatientConditions(id!);
      setConditions(c);
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Failed", "error");
    }
  };

  if (loading) return <p>Loading...</p>;
  if (!patient) return <p>Patient not found</p>;

  return (
    <div>
      <div className="page-header">
        <h1>{patient.first_name} {patient.last_name}</h1>
        <div style={{ display: "flex", gap: 8 }}>
          <Link to={`/app/patients/${id}/edit`} className="btn">Edit</Link>
          <button onClick={handleDelete} className="btn btn-danger">Deactivate</button>
        </div>
      </div>

      <table>
        <tbody>
          <tr><td><strong>Date of Birth</strong></td><td>{patient.date_of_birth || "—"}</td></tr>
          <tr><td><strong>Gender</strong></td><td>{patient.gender || "—"}</td></tr>
          <tr><td><strong>Phone</strong></td><td>{patient.phone || "—"}</td></tr>
          <tr><td><strong>Email</strong></td><td>{patient.email || "—"}</td></tr>
          <tr><td><strong>Address</strong></td><td>{[patient.address, patient.city, patient.state, patient.country].filter(Boolean).join(", ") || "—"}</td></tr>
          <tr><td><strong>Blood Group</strong></td><td>{patient.blood_group || "—"}</td></tr>
          <tr><td><strong>Genotype</strong></td><td>{patient.genotype || "—"}</td></tr>
          <tr><td><strong>Emergency Contact</strong></td><td>{patient.emergency_contact_name ? `${patient.emergency_contact_name} (${patient.emergency_contact_phone})` : "—"}</td></tr>
          <tr><td><strong>Notes</strong></td><td>{patient.notes || "—"}</td></tr>
          <tr><td><strong>Status</strong></td><td><span className={`badge ${patient.is_active ? "badge-active" : "badge-inactive"}`}>{patient.is_active ? "Active" : "Inactive"}</span></td></tr>
        </tbody>
      </table>

      <h2 style={{ marginTop: 32 }}>Allergies</h2>
      <form onSubmit={handleAddAllergy} style={{ display: "flex", gap: 8, marginBottom: 12, alignItems: "end" }}>
        <label style={{ flex: 1 }}>
          Allergy
          <input value={newAllergy} onChange={e => setNewAllergy(e.target.value)} required />
        </label>
        <label>
          Severity
          <select value={newSeverity} onChange={e => setNewSeverity(e.target.value)}>
            <option value="">Select...</option>
            <option value="Mild">Mild</option>
            <option value="Moderate">Moderate</option>
            <option value="Severe">Severe</option>
          </select>
        </label>
        <button type="submit" className="btn">Add</button>
      </form>
      {allergies.length === 0 ? <p>No allergies recorded.</p> : (
        <table>
          <thead><tr><th>Allergy</th><th>Severity</th><th>Notes</th><th>Actions</th></tr></thead>
          <tbody>
            {allergies.map(a => (
              <tr key={a.id}>
                <td>{a.allergy}</td>
                  <td>{a.severity || "—"}</td>
                <td>{a.notes || "—"}</td>
                <td><button onClick={() => handleRemoveAllergy(a.id)} className="action-link action-link-danger">Remove</button></td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      <h2 style={{ marginTop: 32 }}>Conditions</h2>
      <form onSubmit={handleAddCondition} style={{ display: "flex", gap: 8, marginBottom: 12, alignItems: "end" }}>
        <label style={{ flex: 1 }}>
          Condition
          <input value={newCondition} onChange={e => setNewCondition(e.target.value)} required />
        </label>
        <button type="submit" className="btn">Add</button>
      </form>
      {conditions.length === 0 ? <p>No conditions recorded.</p> : (
        <table>
          <thead><tr><th>Condition</th><th>Notes</th></tr></thead>
          <tbody>
            {conditions.map(c => (
              <tr key={c.id}>
                <td>{c.condition}</td>
                <td>{c.notes || "—"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

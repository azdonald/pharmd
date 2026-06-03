import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { getPatient, createPatient, updatePatient } from "../api/patients";

export default function PatientForm() {
  const { id } = useParams();
  const navigate = useNavigate();
  const isNew = !id || id === "new";

  const [form, setForm] = useState({
    first_name: "", last_name: "", date_of_birth: "", gender: "",
    phone: "", email: "", address: "", city: "", state: "", country: "",
    blood_group: "", genotype: "", notes: "",
    emergency_contact_name: "", emergency_contact_phone: "",
  });
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (!isNew && id && id !== "new") {
      getPatient(id).then(p => setForm({
        first_name: p.first_name, last_name: p.last_name, date_of_birth: p.date_of_birth || "",
        gender: p.gender || "", phone: p.phone || "", email: p.email || "",
        address: p.address || "", city: p.city || "", state: p.state || "", country: p.country || "",
        blood_group: p.blood_group || "", genotype: p.genotype || "", notes: p.notes || "",
        emergency_contact_name: p.emergency_contact_name || "", emergency_contact_phone: p.emergency_contact_phone || "",
      })).catch(console.error);
    }
  }, [id, isNew]);

  const update = (field: string) => (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) =>
    setForm({ ...form, [field]: e.target.value });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      if (isNew) {
        const patient = await createPatient(form);
        navigate(`/patients/${patient.id}`);
      } else {
        await updatePatient(id!, form);
        navigate(`/patients/${id}`);
      }
    } catch (err) {
      alert(err instanceof Error ? err.message : "Save failed");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="form-sections-grid">
      <div className="page-header">
        <h1>{isNew ? "New Patient" : "Edit Patient"}</h1>
      </div>

      <form onSubmit={handleSubmit}>
        <div className="form-section">
          <div className="form-section-title">Personal Information</div>
          <div className="form-grid">
            <label>First Name<input value={form.first_name} onChange={update("first_name")} required /></label>
            <label>Last Name<input value={form.last_name} onChange={update("last_name")} required /></label>
            <label>Date of Birth<input type="date" value={form.date_of_birth} onChange={update("date_of_birth")} /></label>
            <label>Gender<select value={form.gender} onChange={update("gender")}>
              <option value="">Select...</option>
              <option value="Male">Male</option>
              <option value="Female">Female</option>
              <option value="Other">Other</option>
            </select></label>
          </div>
        </div>

        <div className="form-section">
          <div className="form-section-title">Contact</div>
          <div className="form-grid">
            <label>Phone<input value={form.phone} onChange={update("phone")} /></label>
            <label>Email<input type="email" value={form.email} onChange={update("email")} /></label>
            <label className="form-grid-full">Address<input value={form.address} onChange={update("address")} /></label>
            <label>City<input value={form.city} onChange={update("city")} /></label>
            <label>State<input value={form.state} onChange={update("state")} /></label>
            <label className="form-grid-full">Country<input value={form.country} onChange={update("country")} /></label>
          </div>
        </div>

        <div className="form-section">
          <div className="form-section-title">Medical</div>
          <div className="form-grid">
            <label>Blood Group<select value={form.blood_group} onChange={update("blood_group")}>
              <option value="">Select...</option>
              <option value="A+">A+</option><option value="A-">A-</option>
              <option value="B+">B+</option><option value="B-">B-</option>
              <option value="AB+">AB+</option><option value="AB-">AB-</option>
              <option value="O+">O+</option><option value="O-">O-</option>
            </select></label>
            <label>Genotype<select value={form.genotype} onChange={update("genotype")}>
              <option value="">Select...</option>
              <option value="AA">AA</option><option value="AS">AS</option><option value="SS">SS</option><option value="AC">AC</option>
            </select></label>
            <label className="form-grid-full">Notes<textarea value={form.notes} onChange={update("notes")} rows={3} /></label>
          </div>
        </div>

        <div className="form-section">
          <div className="form-section-title">Emergency Contact</div>
          <div className="form-grid">
            <label>Name<input value={form.emergency_contact_name} onChange={update("emergency_contact_name")} /></label>
            <label>Phone<input value={form.emergency_contact_phone} onChange={update("emergency_contact_phone")} /></label>
          </div>
        </div>

        <div className="form-actions">
          <button type="submit" className="btn" disabled={saving}>{saving ? "Saving..." : "Save Patient"}</button>
        </div>
      </form>
    </div>
  );
}

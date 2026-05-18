import { useEffect, useState } from "react";
import { listDiscountRules, createDiscountRule, updateDiscountRule, deleteDiscountRule, type DiscountRule } from "../api/pricing";

export default function DiscountRules() {
  const [rules, setRules] = useState<DiscountRule[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);

  const [showForm, setShowForm] = useState(false);
  const [editId, setEditId] = useState("");
  const [formName, setFormName] = useState("");
  const [formType, setFormType] = useState<"percentage" | "fixed">("percentage");
  const [formValue, setFormValue] = useState("");
  const [formAppliesTo, setFormAppliesTo] = useState("all");

  const load = () => {
    setLoading(true);
    listDiscountRules(page, 20).then(res => {
      setRules(res.data);
      setTotal(res.total);
    }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, [page]);

  const totalPages = Math.ceil(total / 20);

  const resetForm = () => {
    setShowForm(false);
    setEditId("");
    setFormName("");
    setFormType("percentage");
    setFormValue("");
    setFormAppliesTo("all");
  };

  const openEdit = (r: DiscountRule) => {
    setEditId(r.id);
    setFormName(r.name);
    setFormType(r.type);
    setFormValue(String(r.value));
    setFormAppliesTo(r.applies_to);
    setShowForm(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formName || !formValue) return;
    try {
      if (editId) {
        await updateDiscountRule(editId, { name: formName, type: formType, value: Number(formValue), applies_to: formAppliesTo as any });
      } else {
        await createDiscountRule({ name: formName, type: formType, value: Number(formValue), applies_to: formAppliesTo as any });
      }
      resetForm();
      load();
    } catch (err) {
      alert(err instanceof Error ? err.message : "Save failed");
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Delete this discount rule?")) return;
    try {
      await deleteDiscountRule(id);
      load();
    } catch (err) {
      alert(err instanceof Error ? err.message : "Delete failed");
    }
  };

  return (
    <div>
      <div className="page-header">
        <h1>Discount Rules</h1>
        {!showForm && <button onClick={() => setShowForm(true)}>New Rule</button>}
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} style={{ marginBottom: 16, padding: 16, border: "1px solid #ddd", borderRadius: 4 }}>
          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr 1fr auto", gap: 8, alignItems: "end" }}>
            <div>
              <label>Name *</label>
              <input value={formName} onChange={e => setFormName(e.target.value)} required />
            </div>
            <div>
              <label>Type</label>
              <select value={formType} onChange={e => setFormType(e.target.value as any)}>
                <option value="percentage">Percentage</option>
                <option value="fixed">Fixed</option>
              </select>
            </div>
            <div>
              <label>Value *</label>
              <input type="number" min={0} step={0.01} value={formValue} onChange={e => setFormValue(e.target.value)} required />
            </div>
            <button type="submit">{editId ? "Update" : "Create"}</button>
            <button type="button" onClick={resetForm}>Cancel</button>
          </div>
          <div style={{ marginTop: 8 }}>
            <label>Applies to</label>
            <select value={formAppliesTo} onChange={e => setFormAppliesTo(e.target.value)}>
              <option value="all">All products</option>
              <option value="category">Category</option>
              <option value="product">Specific product</option>
            </select>
          </div>
        </form>
      )}

      {loading ? <p>Loading...</p> : (
        <>
          <table>
            <thead>
              <tr><th>Name</th><th>Type</th><th>Value</th><th>Applies To</th><th>Active</th><th>Actions</th></tr>
            </thead>
            <tbody>
              {rules.map(r => (
                <tr key={r.id}>
                  <td>{r.name}</td>
                  <td>{r.type}</td>
                  <td>{r.type === "percentage" ? `${r.value}%` : `$${Number(r.value).toFixed(2)}`}</td>
                  <td>{r.applies_to}{r.applies_to_id ? ` (${r.applies_to_id})` : ""}</td>
                  <td>{r.is_active ? "Yes" : "No"}</td>
                  <td>
                    <button onClick={() => openEdit(r)}>Edit</button>
                    <button onClick={() => handleDelete(r.id)}>Delete</button>
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

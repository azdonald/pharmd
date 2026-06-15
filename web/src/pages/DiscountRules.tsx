import { useEffect, useState } from "react";
import { listDiscountRules, createDiscountRule, updateDiscountRule, deleteDiscountRule, type DiscountRule } from "../api/pricing";
import { useToast } from "../context/ToastContext";

function Icon({ name, className }: { name: string; className?: string }) {
  return <span className={`material-symbols-outlined ${className ?? ""}`}>{name}</span>;
}

export default function DiscountRules() {
  const { showToast } = useToast();
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
  const limit = 20;

  const load = () => {
    setLoading(true);
    listDiscountRules(page, limit).then(res => {
      setRules(res.data);
      setTotal(res.total);
    }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, [page]);

  const totalPages = Math.ceil(total / limit);

  const resetForm = () => {
    setShowForm(false);
    setEditId(""); setFormName(""); setFormType("percentage"); setFormValue(""); setFormAppliesTo("all");
  };

  const openEdit = (r: DiscountRule) => {
    setEditId(r.id); setFormName(r.name); setFormType(r.type); setFormValue(String(r.value)); setFormAppliesTo(r.applies_to); setShowForm(true);
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
      showToast(err instanceof Error ? err.message : "Save failed", "error");
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Delete this discount rule?")) return;
    try {
      await deleteDiscountRule(id);
      load();
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Delete failed", "error");
    }
  };

  return (
    <div>
      <div className="mb-8 flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h2 className="font-display-lg text-display-lg text-on-surface">Discount Rules</h2>
          <p className="text-body-lg text-on-surface-variant">Configure pricing discounts and promotions</p>
        </div>
        {!showForm && (
          <button onClick={() => setShowForm(true)}
            className="btn-sky-action">
            <Icon name="add" className="mr-2" />New Rule
          </button>
        )}
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} className="mb-8 p-4 rounded-xl border border-outline-variant bg-surface-container-lowest">
          <h3 className="font-semibold text-on-surface mb-3">{editId ? "Edit Rule" : "New Discount Rule"}</h3>
          <div className="flex flex-col gap-3 sm:flex-row sm:flex-wrap sm:items-end">
            <div>
              <label className="block text-xs font-medium text-on-surface-variant mb-1">Name *</label>
              <input value={formName} onChange={e => setFormName(e.target.value)} required
                className="rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm focus:ring-2 focus:ring-primary outline-none" />
            </div>
            <div>
              <label className="block text-xs font-medium text-on-surface-variant mb-1">Type</label>
              <select value={formType} onChange={e => setFormType(e.target.value as any)}
                className="rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm focus:ring-2 focus:ring-primary outline-none">
                <option value="percentage">Percentage</option>
                <option value="fixed">Fixed</option>
              </select>
            </div>
            <div>
              <label className="block text-xs font-medium text-on-surface-variant mb-1">Value *</label>
              <input type="number" min={0} step={0.01} value={formValue} onChange={e => setFormValue(e.target.value)} required
                className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-primary sm:w-24" />
            </div>
            <div>
              <label className="block text-xs font-medium text-on-surface-variant mb-1">Applies To</label>
              <select value={formAppliesTo} onChange={e => setFormAppliesTo(e.target.value)}
                className="rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-2 text-sm focus:ring-2 focus:ring-primary outline-none">
                <option value="all">All products</option>
                <option value="category">Category</option>
                <option value="product">Specific product</option>
              </select>
            </div>
            <button type="submit"
              className="px-4 py-2 bg-primary text-on-primary font-semibold rounded-lg hover:bg-primary-container transition-all">{editId ? "Update" : "Create"}</button>
            <button type="button" onClick={resetForm}
              className="px-4 py-2 border border-outline-variant rounded-lg text-on-surface hover:bg-surface-container-high transition-all">Cancel</button>
          </div>
        </form>
      )}

      <div className="bg-surface-container-lowest rounded-xl shadow-[0_4px_12px_rgba(0,0,0,0.02)] border border-outline-variant overflow-hidden">
        <div className="overflow-x-auto">
          {loading ? (
            <div className="p-12 text-center"><p className="text-on-surface-variant">Loading discount rules...</p></div>
          ) : rules.length === 0 ? (
            <div className="p-12 text-center"><p className="text-on-surface-variant">No discount rules found</p></div>
          ) : (
            <table className="w-full text-left">
              <thead>
                <tr className="bg-surface-container-low/50">
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Name</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Type</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider text-right">Value</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Applies To</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider text-center">Active</th>
                  <th className="px-6 py-4"></th>
                </tr>
              </thead>
              <tbody className="divide-y divide-outline-variant/30">
                {rules.map(r => (
                  <tr key={r.id} className="hover:bg-surface-container-high/20 transition-colors group">
                    <td className="px-6 py-4 font-semibold text-on-surface">{r.name}</td>
                    <td className="px-6 py-4"><code className="text-sm bg-surface-container-high px-2 py-0.5 rounded">{r.type}</code></td>
                    <td className="px-6 py-4 text-right font-data-mono text-data-mono">
                      {r.type === "percentage" ? `${r.value}%` : `$${Number(r.value).toFixed(2)}`}
                    </td>
                    <td className="px-6 py-4 text-body-md text-on-surface-variant">
                      {r.applies_to}{r.applies_to_id ? ` (${r.applies_to_id})` : ""}
                    </td>
                    <td className="px-6 py-4 text-center">
                      <span className={`inline-flex items-center px-3 py-1 rounded-full text-label-caps font-bold ${
                        r.is_active ? "bg-secondary-container/20 text-secondary" : "bg-outline-variant/20 text-on-surface-variant"
                      }`}>{r.is_active ? "Yes" : "No"}</span>
                    </td>
                    <td className="px-6 py-4 text-right">
                      <div className="flex items-center justify-end space-x-2">
                        <button onClick={() => openEdit(r)} className="px-3 py-1 text-sm text-primary hover:bg-primary/5 rounded transition-colors">Edit</button>
                        <button onClick={() => handleDelete(r.id)} className="px-3 py-1 text-sm text-error hover:bg-error/5 rounded transition-colors">Delete</button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
        {totalPages > 1 && (
          <div className="px-6 py-4 border-t border-outline-variant flex justify-between items-center bg-surface-container-low/10">
            <p className="text-body-md text-on-surface-variant">Showing {((page - 1) * limit) + 1} to {Math.min(page * limit, total)} of {total} rules</p>
            <div className="flex space-x-2">
              <button onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page === 1}
                className="p-2 border border-outline-variant rounded-md hover:bg-surface-container-high disabled:opacity-50 disabled:cursor-not-allowed"><Icon name="chevron_left" /></button>
              {Array.from({ length: Math.min(totalPages, 5) }, (_, i) => {
                const start = Math.max(1, Math.min(page - 2, totalPages - 4)); const p = start + i;
                if (p > totalPages) return null;
                return (
                  <button key={p} onClick={() => setPage(p)}
                    className={`w-10 h-10 rounded-md flex items-center justify-center font-bold text-sm ${
                      p === page ? "bg-primary text-on-primary" : "border border-outline-variant hover:bg-surface-container-high"
                    }`}>{p}</button>
                );
              })}
              <button onClick={() => setPage(p => Math.min(totalPages, p + 1))} disabled={page === totalPages}
                className="p-2 border border-outline-variant rounded-md hover:bg-surface-container-high disabled:opacity-50 disabled:cursor-not-allowed"><Icon name="chevron_right" /></button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

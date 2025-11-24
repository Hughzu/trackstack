# The Decision Framework

## Use Astro Components When:

- Content is read-only or rarely changes
- User interaction is simple (clicks, basic navigation)
- Data is fetched once at build/request time
- No complex state management needed

## Use Angular Components When:

- Forms with validation and complex interactions
- Real-time updates (charts, counters, live data)
- Client-side filtering/sorting of large datasets
- Multi-step workflows with state
- Complex UI interactions (drag-and-drop, rich editors)

## Rule of Thumb

- If it can be a `<form>` with a submit button → **Use Astro**
- If it needs instant feedback while typing → **Use Angular**

---

## Business Capability Examples

### Expense Tracker

#### ✅ Use Astro: Expense List (Read-only)

```astro
---
// src/modules/expense-tracker/components/ExpenseList.astro
import { expenseRepository } from '../repositories/expenseRepository';

const userId = Astro.locals.user.id;
const expenses = expenseRepository.findByUserId(userId);
---

<div class="expense-list">
  <h2>Recent Expenses</h2>
  {expenses.map(expense => (
    <div class="expense-item">
      <span class="amount">€{expense.amount.toFixed(2)}</span>
      <span class="category">{expense.category}</span>
      <span class="date">{new Date(expense.date).toLocaleDateString()}</span>
      <a href={`/expenses/${expense.id}`}>View</a>
    </div>
  ))}
</div>
```

**Why Astro?**
- Static list of data
- No client-side interaction beyond navigation
- Fast initial render
- Zero JavaScript shipped

---

#### ✅ Expense Detail Page (Mostly Static)

```astro
---
// src/pages/expenses/[id].astro
import AppLayout from '@/shared/layouts/AppLayout.astro';
import { expenseRepository } from '@/modules/expense-tracker/repositories/expenseRepository';

const { id } = Astro.params;
const expense = expenseRepository.findById(id);

if (!expense) {
  return Astro.redirect('/expenses');
}
---

<AppLayout title="Expense Details">
  <h1>Expense Details</h1>
  
  <dl>
    <dt>Amount</dt>
    <dd>€{expense.amount.toFixed(2)}</dd>
    
    <dt>Category</dt>
    <dd>{expense.category}</dd>
    
    <dt>Description</dt>
    <dd>{expense.description}</dd>
    
    <dt>Date</dt>
    <dd>{new Date(expense.date).toLocaleDateString()}</dd>
  </dl>

  <!-- Simple form actions with Astro -->
  <div class="actions">
    <a href={`/expenses/${id}/edit`}>Edit</a>
    <form method="POST" action={`/api/expenses/${id}/delete`}>
      <button type="submit">Delete</button>
    </form>
  </div>
</AppLayout>
```

**Why Astro?**
- Read-only content display
- Delete action uses native form submission (no JS needed)
- Edit links to a different page

---

#### ⚡ Use Angular: Expense Form (Create/Edit with Validation)

```typescript
// src/modules/expense-tracker/components/ExpenseForm.component.ts
import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'expense-form',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <form (ngSubmit)="onSubmit()" class="expense-form">
      <div class="form-group">
        <label for="amount">Amount *</label>
        <input 
          type="number" 
          id="amount"
          [(ngModel)]="expense.amount" 
          name="amount"
          (input)="validateAmount()"
          [class.error]="errors.amount"
          required
        />
        <span class="error-message" *ngIf="errors.amount">
          {{ errors.amount }}
        </span>
      </div>
      
      <div class="form-group">
        <label for="category">Category *</label>
        <select 
          id="category"
          [(ngModel)]="expense.category" 
          name="category"
          (change)="onCategoryChange()"
        >
          <option value="">Select category</option>
          <option value="food">🍔 Food</option>
          <option value="transport">🚗 Transport</option>
          <option value="utilities">💡 Utilities</option>
          <option value="entertainment">🎮 Entertainment</option>
        </select>
      </div>
      
      <!-- Conditional field based on category -->
      <div class="form-group" *ngIf="expense.category === 'transport'">
        <label for="distance">Distance (km)</label>
        <input 
          type="number" 
          id="distance"
          [(ngModel)]="expense.distance" 
          name="distance"
        />
        <p class="helper-text">Cost per km: €{{ costPerKm }}</p>
      </div>
      
      <div class="form-group">
        <label for="description">Description</label>
        <textarea 
          id="description"
          [(ngModel)]="expense.description" 
          name="description"
          rows="3"
          (input)="updateCharCount()"
        ></textarea>
        <span class="char-count">{{ charCount }}/500</span>
      </div>
      
      <div class="form-group">
        <label for="date">Date *</label>
        <input 
          type="date" 
          id="date"
          [(ngModel)]="expense.date" 
          name="date"
          [max]="today"
          required
        />
      </div>
      
      <!-- Live calculation -->
      <div class="summary">
        <strong>Total: €{{ expense.amount || 0 }}</strong>
        <span class="budget-indicator" [class.over-budget]="isOverBudget()">
          {{ getBudgetStatus() }}
        </span>
      </div>
      
      <button 
        type="submit" 
        [disabled]="!isValid()"
        class="submit-btn"
      >
        {{ isEditMode ? 'Update' : 'Create' }} Expense
      </button>
    </form>
  `
})
export class ExpenseFormComponent {
  @Input() initialData?: any;
  @Input() isEditMode = false;
  @Output() submitted = new EventEmitter();

  expense = {
    amount: 0,
    category: '',
    description: '',
    date: new Date().toISOString().split('T')[0],
    distance: 0
  };

  errors: any = {};
  charCount = 0;
  today = new Date().toISOString().split('T')[0];
  costPerKm = 0.35;
  monthlyBudget = 500;

  ngOnInit() {
    if (this.initialData) {
      this.expense = { ...this.initialData };
    }
  }

  validateAmount() {
    if (this.expense.amount <= 0) {
      this.errors.amount = 'Amount must be greater than 0';
    } else if (this.expense.amount > 10000) {
      this.errors.amount = 'Amount seems too high. Please verify.';
    } else {
      this.errors.amount = null;
    }
  }

  onCategoryChange() {
    // Reset category-specific fields
    if (this.expense.category !== 'transport') {
      this.expense.distance = 0;
    }
  }

  updateCharCount() {
    this.charCount = this.expense.description?.length || 0;
  }

  isOverBudget(): boolean {
    // Could fetch monthly total from API
    return this.expense.amount > 100; // Simplified
  }

  getBudgetStatus(): string {
    const remaining = this.monthlyBudget - this.expense.amount;
    return remaining > 0 
      ? `€${remaining} remaining this month`
      : `€${Math.abs(remaining)} over budget!`;
  }

  isValid(): boolean {
    return this.expense.amount > 0 
      && this.expense.category !== '' 
      && this.expense.date !== ''
      && !this.errors.amount;
  }

  onSubmit() {
    if (this.isValid()) {
      this.submitted.emit(this.expense);
    }
  }
}
```

**Why Angular?**
- Real-time validation as user types
- Dynamic fields (distance only shows for transport)
- Live calculations (budget status, char count)
- Complex state (multiple form fields with interdependencies)
- Instant feedback improves UX significantly

---

#### Usage in Astro page:

```astro
---
// src/pages/expenses/new.astro
import AppLayout from '@/shared/layouts/AppLayout.astro';
import { ExpenseFormComponent } from '@/modules/expense-tracker';
---

<AppLayout title="Add New Expense">
  <h1>Add New Expense</h1>
  <expense-form client:load />
</AppLayout>
```

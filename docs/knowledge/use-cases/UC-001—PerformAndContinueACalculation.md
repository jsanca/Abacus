UC-001 — Perform and Continue a Calculation
Metadata
Field	Value
ID	UC-001
Name	Perform and Continue a Calculation
Status	Approved
Primary actor	Calculator user
Trigger	The user opens Abacus to perform an arithmetic operation
Scope	React client and calculator API
Goal

Allow a user to select an available arithmetic operation, provide the required operands, obtain a result from the backend, and optionally continue calculating from that result.

Preconditions
The Abacus client is available.
The calculator backend is reachable.
The client has successfully loaded the operation manifest.
At least one operation is available.
Supported operations
Operation	Arity	Presentation	Expected behavior
Addition	2	Infix a + b	Returns the sum
Subtraction	2	Infix a − b	Returns the difference
Multiplication	2	Infix a × b	Returns the product
Division	2	Infix a ÷ b	Returns the quotient when validation succeeds
Exponentiation	2	Infix a xʸ b or equivalent accessible notation	Raises the first operand to the second
Square Root	1	Prefix √a	Returns the real square root when validation succeeds
Percentage	2	b% of a	Returns b percent of the first operand

For percentage, the second operand accepts values from 0 through 100, inclusive.

Main success flow
The client requests the operation manifest from the backend.
The client displays all available operations as directly selectable controls.
The user selects an operation with one click or keyboard activation.
The client renders one or two operand inputs according to the operation arity.
The user enters the required operand values.
The user activates Calculate or presses Enter.
The client performs structural input validation.
The client evaluates the operation validations projected by the backend manifest.
The client sends the operation identifier and operands to the backend.
The backend independently validates the request.
The backend executes the registered operation.
The client displays the completed expression and result.
The result is announced accessibly.
The completed expression remains visible until the user begins another calculation.
Continue-calculation flow

Given:

144 ÷ 12 = 12

When the user selects multiplication:

The result 12 becomes Operand 1.
Multiplication becomes the selected operation.
Operand 2 is cleared.
The previous resolved expression transitions into an editable calculation.
Focus moves to Operand 2.

The user can then enter 3 and calculate:

12 × 3 = 36
Alternative flows
A1 — Unary operation
The user selects square root.
The client shows only Operand 1.
Operand 2 is absent from the layout and tab order.
The expression is displayed using prefix notation.
The normal calculation flow continues.
A2 — Percentage
The user selects percentage.
Operand 1 represents the base value.
Operand 2 represents the percentage.
Operand 2 displays a fixed, non-editable % suffix.
The client evaluates manifest-provided range rules.
The result is presented as percentage% of base.
A3 — User changes the operator while editing
The existing operand values are preserved where compatible.
Arity and operand presentation update to the newly selected operation.
Stale validation errors are cleared or reevaluated.
Focus moves only when doing so improves the editing flow.
A4 — User clears the calculator
The user activates Clear or presses Escape.
Operand values are cleared.
Result and completed expression are cleared.
Validation and service errors are cleared.
The initial or default operation is selected.
Focus returns to Operand 1.
Validation flows
V1 — Missing required operand

The client prevents submission and identifies the missing operand.

V2 — Structurally invalid numeric input

The client prevents submission when an operand cannot be represented as a supported number.

V3 — Manifest validation failure

The frontend generically evaluates a validation supplied by the backend manifest.

The frontend must not contain operation-specific branches such as:

operation.id === "division"

to enforce mathematical rules.

V4 — Backend validation failure

When the backend rejects a request, the client displays the returned violation near the affected operand when targeting information is available.

Failure flows
F1 — Operation manifest unavailable
The calculator controls are not rendered from hardcoded fallback data.
The user sees a service-unavailable state.
A Retry action is available.
F2 — Calculation request fails
The current operands and selected operation remain intact.
A service error is displayed.
The user can retry without re-entering the calculation.
F3 — Unknown operation
The client displays a stable application error.
No incorrect result is shown.
Postconditions
Success
A backend-produced result is visible.
The completed expression remains available for context.
The result can become the first operand of another operation.
Failure
No invalid result is displayed.
User-entered values are preserved unless Clear was explicitly requested.
The user receives actionable feedback.
Business and interaction rules
ID	Rule
BR-001	Operations are discovered from the backend manifest.
BR-002	Every operation declares an arity of one or two.
BR-003	The frontend performs only structural validation and generic interpretation of projected validation expressions.
BR-004	The backend remains authoritative and independently validates every request.
BR-005	All seven required operations must be visible and selectable with one action.
BR-006	Division by zero must be rejected by the validation model.
BR-007	Square root accepts only values allowed by its projected real-number validation.
BR-008	Percentage accepts values from 0 through 100 inclusive for Operand 2.
BR-009	Selecting an operator after a successful calculation continues from the previous result.
BR-010	The completed expression remains visible until a new calculation begins.
BR-011	Enter calculates when submission is valid.
BR-012	Escape clears the calculator.
BR-013	Backend and service errors must not erase user input.
BR-014	The application must be usable at desktop and mobile viewport sizes.
BR-015	The result and errors must be accessible to assistive technologies.

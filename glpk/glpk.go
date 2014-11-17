// This code is part of glpk package (Go bindings for the GNU Linear Programming Kit).
//
// Copyright (C) 2014 Łukasz Pankowski <lukpank@o2.pl>
//
// Some comments/strings are taken or adapted from GLPK and thus are
// subject to the following copyright:
//
// Copyright (C) 2000, 2001, 2002, 2003, 2004, 2005, 2006, 2007, 2008,
// 2009, 2010, 2011, 2013, 2014 Andrew Makhorin, Department for Applied
// Informatics, Moscow Aviation Institute, Moscow, Russia. All rights
// reserved. E-mail: <mao@gnu.org>.
//
// Pacakge glpk is free software: you can redistribute it and/or
// modify it under the terms of the GNU General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// Package glpk is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with glpk package. If not, see <http://www.gnu.org/licenses/>.

// Go bindings for GLPK (GNU Linear Programming Kit).
//
// For a usage example see https://github.com/lukpank/go-glpk#example.
//
// The binding is not complete but enough for my purposes. Fill free
// to contact me if there is some part of GLPK that you would like to
// use and it is not yet covered by the glpk package.
//
// Package glpk is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
package glpk

import (
	"reflect"
	"runtime"
	"unsafe"
)

// #cgo LDFLAGS: -lglpk
// #include <glpk.h>
// #include <stdlib.h>
import "C"

// Objective function direction (maximization or minimization).
type ObjDir int

const (
	MAX = ObjDir(C.GLP_MAX) // MAX represents maximization
	MIN = ObjDir(C.GLP_MIN) // MIN represents minimization
)

// Bounds type of a variable
type BndsType int

const (
	FR = BndsType(C.GLP_FR) // FR represents a free (unbounded) variable
	LO = BndsType(C.GLP_LO) // LO represents a lower-bounded variable
	UP = BndsType(C.GLP_UP) // UP represents an upper-bounded variable
	DB = BndsType(C.GLP_DB) // DB represents a double-bounded variable
	FX = BndsType(C.GLP_FX) // FX represents a fixed variable
)

// Solution Status
type SolStat int

const (
	UNDEF  = SolStat(C.GLP_UNDEF)  // UNDEF indicates that solution is undefined
	FEAS   = SolStat(C.GLP_FEAS)   // FEAS indicates that solution is feasible
	INFEAS = SolStat(C.GLP_INFEAS) // INFEAS indicates that the solution is infeasible
	NOFEAS = SolStat(C.GLP_NOFEAS) // NOFEAS indicates that there is no feasible solution
	OPT    = SolStat(C.GLP_OPT)    // OPT indicates that the solution is optimal
	UNBND  = SolStat(C.GLP_UNBND)  // UNBND indicates that the problem has unbounded solution
)

type prob struct {
	p *C.glp_prob
}

// Prob represens optimization problem. Use glpk.New() to create a new problem.
type Prob struct {
	p *prob
}

func finalizeProb(p *prob) {
	if p.p != nil {
		C.glp_delete_prob(p.p)
		p.p = nil
	}
}

// New creates a new optimization problem.
func New() *Prob {
	p := &prob{C.glp_create_prob()}
	runtime.SetFinalizer(p, finalizeProb)
	return &Prob{p}
}

// Delete deletes a problem.  Calling Delete on a deleted problem will
// have no effect (It is save to do so). But calling any other method
// on a deleted problem will panic. The problem will be deleted on
// garbage collection but you can do this as soon as you no longer
// need the optimization problem.
func (p *Prob) Delete() {
	if p.p.p != nil {
		C.glp_delete_prob(p.p.p)
		p.p.p = nil
	}
}

// Erase erases the problem. After erasing the problem is empty as if
// it were created with glpk.New().
func (p *Prob) Erase() {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	C.glp_erase_prob(p.p.p)
}

// SetProbName sets (changes) the problem name.
func (p *Prob) SetProbName(name string) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	C.glp_set_prob_name(p.p.p, s)
}

// SetObjName sets (changes) objective function name.
func (p *Prob) SetObjName(name string) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	C.glp_set_obj_name(p.p.p, s)
}

// SetObjDir sets optimization direction (either glpk.MAX for
// maximization or glpk.MIN for minimization)
func (p *Prob) SetObjDir(dir ObjDir) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	C.glp_set_obj_dir(p.p.p, C.int(dir))
}

// AddRows adds rows (constraints). Returns (1-based) index of the
// first of the added rows.
func (p *Prob) AddRows(nrs int) int {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return int(C.glp_add_rows(p.p.p, C.int(nrs)))
}

// AddCols adds columns (variables). Returns (1-based) index of the
// first of the added columns.
func (p *Prob) AddCols(nrs int) int {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return int(C.glp_add_cols(p.p.p, C.int(nrs)))
}

// SetRowName sets i-th row (constraint) name.
func (p *Prob) SetRowName(i int, name string) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	C.glp_set_row_name(p.p.p, C.int(i), s)
}

// SetColName sets j-th column (variable) name.
func (p *Prob) SetColName(j int, name string) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	C.glp_set_col_name(p.p.p, C.int(j), s)
}

// SetRowBnds sets row bounds
func (p *Prob) SetRowBnds(i int, type_ BndsType, lb float64, ub float64) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	C.glp_set_row_bnds(p.p.p, C.int(i), C.int(type_), C.double(lb), C.double(ub))
}

// SetColBnds sets column bounds
func (p *Prob) SetColBnds(j int, type_ BndsType, lb float64, ub float64) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	C.glp_set_col_bnds(p.p.p, C.int(j), C.int(type_), C.double(lb), C.double(ub))
}

// SetObjCoef sets objective function coefficient of j-th column.
func (p *Prob) SetObjCoef(j int, coef float64) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	C.glp_set_obj_coef(p.p.p, C.int(j), C.double(coef))
}

// SetMatRow sets (replaces) i-th row. It sets
//
//     matrix[i, ind[j]] = val[j]
//
// for j=1..len(ind). ind[0] and val[0] are ignored. Requires
// len(ind) = len(val).
func (p *Prob) SetMatRow(i int, ind []int32, val []float64) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	if len(ind) != len(val) {
		panic("len(ind) and len(val) should be equal")
	}
	ind_ := (*reflect.SliceHeader)(unsafe.Pointer(&ind))
	val_ := (*reflect.SliceHeader)(unsafe.Pointer(&val))
	C.glp_set_mat_row(p.p.p, C.int(i), C.int(len(ind)-1), (*C.int)(unsafe.Pointer(ind_.Data)), (*C.double)(unsafe.Pointer(val_.Data)))
}

// SetMatCol sets (replaces) j-th column. It sets
//
//     matrix[ind[i], j] = val[i]
//
// for i=1..len(ind). ind[0] and val[0] are ignored. Requires
// len(ind) = len(val).
func (p *Prob) SetMatCol(j int, ind []int32, val []float64) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	if len(ind) != len(val) {
		panic("len(ind) and len(val) should be equal")
	}
	ind_ := (*reflect.SliceHeader)(unsafe.Pointer(&ind))
	val_ := (*reflect.SliceHeader)(unsafe.Pointer(&val))
	C.glp_set_mat_col(p.p.p, C.int(j), C.int(len(ind)-1), (*C.int)(unsafe.Pointer(ind_.Data)), (*C.double)(unsafe.Pointer(val_.Data)))
}

// LoadMatrix replaces all of the constraint matrix. It sets
//
//     matrix[ia[i], ja[i]] = ar[i]
//
// for i = 1..len(ia). ia[0], ja[0], and ar[0] are ignored. It
// requiers len(ia)=len(ja)=len(ar).
func (p *Prob) LoadMatrix(ia, ja []int32, ar []float64) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	if len(ia) != len(ja) || len(ia) != len(ar) {
		panic("len(ia) and len(ja) and len(ar) should be equal")
	}
	ia_ := (*reflect.SliceHeader)(unsafe.Pointer(&ia))
	ja_ := (*reflect.SliceHeader)(unsafe.Pointer(&ja))
	ar_ := (*reflect.SliceHeader)(unsafe.Pointer(&ar))
	C.glp_load_matrix(p.p.p, C.int(len(ia)-1), (*C.int)(unsafe.Pointer(ia_.Data)), (*C.int)(unsafe.Pointer(ja_.Data)), (*C.double)(unsafe.Pointer(ar_.Data)))
}

// TODO:
// glp_check_dup
// glp_del_rows

// Copy returns a copy of the given optimization problem. If name is
// true also symbolic names are copies otherwise their not copied
func (p *Prob) Copy(names bool) *Prob {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	q := &Prob{&prob{C.glp_create_prob()}}
	var names_ C.int
	if names {
		names_ = C.GLP_ON
	} else {
		names_ = C.GLP_OFF
	}
	C.glp_copy_prob(q.p.p, p.p.p, names_)
	return q
}

// ProbName returns problem name.
func (p *Prob) ProbName() string {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return C.GoString(C.glp_get_prob_name(p.p.p))
}

// ObjName returns objective name.
func (p *Prob) ObjName() string {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return C.GoString(C.glp_get_obj_name(p.p.p))
}

// ObjDir returns optimization direction (either glpk.MAX or glpk.MIN).
func (p *Prob) ObjDir() ObjDir {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return ObjDir(C.glp_get_obj_dir(p.p.p))
}

// NumRows returns number of rows.
func (p *Prob) NumRows() int {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return int(C.glp_get_num_rows(p.p.p))
}

// NumCols returns number of columns.
func (p *Prob) NumCols() int {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return int(C.glp_get_num_cols(p.p.p))
}

// RowName returns row (constraint) name of i-th row.
func (p *Prob) RowName(i int) string {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return C.GoString(C.glp_get_row_name(p.p.p, C.int(i)))
}

// ColName returns column (variable) name of j-th column.
func (p *Prob) ColName(j int) string {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return C.GoString(C.glp_get_col_name(p.p.p, C.int(j)))
}

// TODO:
// glp_get_row_type
// glp_get_row_lb
// glp_get_row_ub
// glp_get_col_type
// glp_get_col_lb
// glp_get_col_ub

// ObjCoef returns objective function coefficient of j-th column.
func (p *Prob) ObjCoef(j int) float64 {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return float64(C.glp_get_obj_coef(p.p.p, C.int(j)))
}

// TODO:
// glp_get_num_nz

// MatRow returns nonzero elements of i-th row. ind[1]..ind[n] are
// column numbers of the nonzero elements of the row, val[1]..val[n]
// are their values, and n is the number of nonzero elements in the
// row.
func (p *Prob) MatRow(i int) (ind []int32, val []float64) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	if len(ind) != len(val) {
		panic("len(ind) and len(val) should be equal")
	}
	length := C.glp_get_mat_row(p.p.p, C.int(i), nil, nil)
	ind = make([]int32, length+1)
	val = make([]float64, length+1)
	ind_ := (*reflect.SliceHeader)(unsafe.Pointer(&ind))
	val_ := (*reflect.SliceHeader)(unsafe.Pointer(&val))
	C.glp_get_mat_row(p.p.p, C.int(i), (*C.int)(unsafe.Pointer(ind_.Data)), (*C.double)(unsafe.Pointer(val_.Data)))
	return
}

// MatCol returns nonzero elements of j-th column. ind[1]..ind[n] are
// row numbers of the nonzero elements of the column, val[1]..val[n]
// are their values, and n is the number of nonzero elements in the
// column.
func (p *Prob) MatCol(j int) (ind []int32, val []float64) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	if len(ind) != len(val) {
		panic("len(ind) and len(val) should be equal")
	}
	length := C.glp_get_mat_col(p.p.p, C.int(j), nil, nil)
	ind = make([]int32, length+1)
	val = make([]float64, length+1)
	ind_ := (*reflect.SliceHeader)(unsafe.Pointer(&ind))
	val_ := (*reflect.SliceHeader)(unsafe.Pointer(&val))
	C.glp_get_mat_col(p.p.p, C.int(j), (*C.int)(unsafe.Pointer(ind_.Data)), (*C.double)(unsafe.Pointer(val_.Data)))
	return
}

// TODO:
// glp_create_index
// glp_find_row
// glp_find_col
// glp_delete_index
// glp_set_rii
// glp_set_sjj
// glp_get_rii
// glp_get_sjj
// glp_scale_prob
// glp_unscale_prob
// glp_set_row_stat
// glp_set_col_stat
// glp_std_basis
// glp_adv_basis
// glp_cpx_basis

// Optimization Error
type OptError int

const (
	EBADB   = OptError(C.GLP_EBADB)   // invalid basis
	ESING   = OptError(C.GLP_ESING)   // singular matrix
	ECOND   = OptError(C.GLP_ECOND)   // ill-conditioned matrix
	EBOUND  = OptError(C.GLP_EBOUND)  // invalid bounds
	EFAIL   = OptError(C.GLP_EFAIL)   // solver failed
	EOBJLL  = OptError(C.GLP_EOBJLL)  // objective lower limit reached
	EOBJUL  = OptError(C.GLP_EOBJUL)  // objective upper limit reached
	EITLIM  = OptError(C.GLP_EITLIM)  // iteration limit exceeded
	ETMLIM  = OptError(C.GLP_ETMLIM)  // time limit exceeded
	ENOPFS  = OptError(C.GLP_ENOPFS)  // no primal feasible solution
	ENODFS  = OptError(C.GLP_ENODFS)  // no dual feasible solution
	EROOT   = OptError(C.GLP_EROOT)   // root LP optimum not provided
	ESTOP   = OptError(C.GLP_ESTOP)   // search terminated by application
	EMIPGAP = OptError(C.GLP_EMIPGAP) // relative mip gap tolerance reached
	ENOFEAS = OptError(C.GLP_ENOFEAS) // no primal/dual feasible solution
	ENOCVG  = OptError(C.GLP_ENOCVG)  // no convergence
	EINSTAB = OptError(C.GLP_EINSTAB) // numerical instability
	EDATA   = OptError(C.GLP_EDATA)   // invalid data
	ERANGE  = OptError(C.GLP_ERANGE)  // result out of range
)

func (r OptError) Error() string {
	switch r {
	case EBADB:
		return "invalid basis"
	case ESING:
		return "singular matrix"
	case ECOND:
		return "ill-conditioned matrix"
	case EBOUND:
		return "invalid bounds"
	case EFAIL:
		return "solver failed"
	case EOBJLL:
		return "objective lower limit reached"
	case EOBJUL:
		return "objective upper limit reached"
	case EITLIM:
		return "iteration limit exceeded"
	case ETMLIM:
		return "time limit exceeded"
	case ENOPFS:
		return "no primal feasible solution"
	case ENODFS:
		return "no dual feasible solution"
	case EROOT:
		return "root LP optimum not provided"
	case ESTOP:
		return "search terminated by application"
	case EMIPGAP:
		return "relative mip gap tolerance reached"
	case ENOFEAS:
		return "no primal/dual feasible solution"
	case ENOCVG:
		return "no convergence"
	case EINSTAB:
		return "numerical instability"
	case EDATA:
		return "invalid data"
	case ERANGE:
		return "result out of range"
	}
	return "unknown error"
}

// Simplex solves LP with Simplex method. The argument parm may by nil
// (means that default values will be used). See also NewSmcp().
// Returns nil if problem have been solved (not necessarly finding
// optimal solution) otherwise returns an error which is an instanse
// of OptError.
func (p *Prob) Simplex(parm *Smcp) error {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	var err OptError
	if parm != nil {
		err = OptError(C.glp_simplex(p.p.p, &parm.smcp))
	} else {
		err = OptError(C.glp_simplex(p.p.p, nil))
	}
	if err == 0 {
		return nil
	}
	return err
}

// Exact solves LP with Simplex method using exact (rational)
// arithmetic. argument parm may by nil (means that default values
// will be used). See also NewSmcp().  Returns nil if problem have
// been solved (not necessarly finding optimal solution) otherwise
// returns an error which is an instanse of OptError.
func (p *Prob) Exact(parm *Smcp) error {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	var err OptError
	if parm != nil {
		err = OptError(C.glp_exact(p.p.p, &parm.smcp))
	} else {
		err = OptError(C.glp_exact(p.p.p, nil))
	}
	if err == 0 {
		return nil
	}
	return err
}

// Smcp represents simplex solver control parameters, a set of
// parameters for Prob.Simplex() and Prob.Exact(). Please use
// NewSmcp() to create Smtp structure which is properly initialized.
type Smcp struct {
	smcp C.glp_smcp
}

// NewSmcp creates new Smcp struct (a set of simplex solver control
// parameters) to be given as argument of Prob.Simplex() or
// Prob.Exact().
func NewSmcp() *Smcp {
	s := new(Smcp)
	C.glp_init_smcp(&s.smcp)
	return s
}

// Message level
type MsgLev int

const (
	// Message levels (default: glpk.MSG_ALL). Usage example:
	//
	//     lp := glpk.New()
	//     ...
	//     smcp := glpk.NewSmcp()
	//     smcp.SetMsgLev(glpk.MSG_ERR)
	//     lp.Simplex(smcp)
	//
	MSG_OFF = MsgLev(C.GLP_MSG_OFF) // no output
	MSG_ERR = MsgLev(C.GLP_MSG_ERR) // warning and error messages only
	MSG_ON  = MsgLev(C.GLP_MSG_ON)  // normal output
	MSG_ALL = MsgLev(C.GLP_MSG_ALL) // full output
	MSG_DBG = MsgLev(C.GLP_MSG_DBG) // debug output
)

// SetMsgLev sets message level displayed by the optimization function
// (default: glpk.MSG_ALL).
func (s *Smcp) SetMsgLev(lev MsgLev) {
	s.smcp.msg_lev = C.int(lev)
}

// Simplex method option
type Meth int

const (
	// Simplex method options (default: glpk.PRIMAL). Usage example:
	//
	//     lp := glpk.New()
	//     ...
	//     smcp := glpk.NewSmcp()
	//     smcp.SetMeth(glpk.DUALP)
	//     lp.Simplex(smcp)
	//
	PRIMAL = Meth(C.GLP_PRIMAL) // use primal simplex
	DUALP  = Meth(C.GLP_DUALP)  // use dual; if it fails, use primal
	DUAL   = Meth(C.GLP_DUAL)   // use dual simplex
)

// SetMeth sets simplex method option (default: glpk.PRIMAL).
func (s *Smcp) SetMeth(meth Meth) {
	s.smcp.meth = C.int(meth)
}

// Pricing technique
type Pricing int

const (
	// Pricing techniques (default: glpk.PT_PSE). Example usage
	//
	//     lp := glpk.New()
	//     ...
	//     smcp := glpk.NewSmcp()
	//     smcp.SetPricing(glpk.PT_STD)
	//     lp.Simplex(smcp)
	//
	PT_STD = Pricing(C.GLP_PT_STD) // standard (Dantzig rule)
	PT_PSE = Pricing(C.GLP_PT_PSE) // projected steepest edge
)

// SetPricing sets pricing technique (default: glpk.PT_PSE).
func (s *Smcp) SetPricing(pricing Pricing) {
	s.smcp.pricing = C.int(pricing)
}

// Ratio test technique
type RTest int

const (
	// Ratio test techniques (default: glpk.RT_HAR). Example usage:
	//
	//     lp := glpk.New()
	//     ...
	//     smcp := glpk.NewSmcp()
	//     smcp.SetRTest(glpk.RT_STD)
	//     lp.Simplex(smcp)
	//
	RT_STD = RTest(C.GLP_RT_STD) // standard (textbook)
	RT_HAR = RTest(C.GLP_RT_HAR) // two-pass Harris' ratio test
)

// SetRTest sets ratio test technique (default: glpk.RT_HAR)
func (s *Smcp) SetRTest(r_test RTest) {
	s.smcp.r_test = C.int(r_test)
}

// Status returns status of the basic solution.
func (p *Prob) Status() SolStat {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return SolStat(C.glp_get_status(p.p.p))
}

// PrimStat returns status of the primal basic solution.
func (p *Prob) PrimStat() SolStat {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return SolStat(C.glp_get_prim_stat(p.p.p))
}

// DualStat returns status of the dual basic solution.
func (p *Prob) DualStat() SolStat {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return SolStat(C.glp_get_dual_stat(p.p.p))
}

// ObjVal returns objective function value.
func (p *Prob) ObjVal() float64 {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return float64(C.glp_get_obj_val(p.p.p))
}

// TODO:
// glp_get_row_stat
// glp_get_row_prim
// glp_get_row_dual
// glp_get_col_stat

// ColPrim returns primal value of the variable associated with j-th
// column.
func (p *Prob) ColPrim(j int) float64 {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return float64(C.glp_get_col_prim(p.p.p, C.int(j)))
}

// TODO:
// glp_get_col_dual
// ...

type tran struct {
	t *C.glp_tran
}

type Tran struct {
	t *tran
}

func finalizeTran(t *tran) {
	if t.t != nil {
		C.glp_mpl_free_wksp(t.t)
		t.t = nil
	}
}

func NewMpl() *Tran {
	t := &tran{C.glp_mpl_alloc_wksp()}
	runtime.SetFinalizer(t, finalizeTran)
	return &Tran{t}
}

func (t *Tran) MplFreeWksp() {
	if t.t.t != nil {
		C.glp_mpl_free_wksp(t.t.t)
		t.t.t = nil
	}
}

func (t *Tran) MplReadModel(filename string, skipDataFlag bool) int {
	f := C.CString(filename)

	skip := C.int(0)
	if skipDataFlag == true {
		skip = C.int(1)
	}

	ret := C.glp_mpl_read_model(t.t.t, f, skip)

	return int(ret)
}

func (t *Tran) MplGenerate() int {

	ret := C.glp_mpl_generate(t.t.t, nil)

	return int(ret)
}

func (t *Tran) MplBuildProb(p *Prob) {
	C.glp_mpl_build_prob(t.t.t, p.p.p)
}

func (t *Tran) MplReadData(filename string) int {
	f := C.CString(filename)

	ret := C.glp_mpl_read_data(t.t.t, f)

	return int(ret)
}

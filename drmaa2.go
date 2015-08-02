/*
Copyright 2014, 2015 Daniel Gruber, http://www.gridengine.eu

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

// Package drmaa2 implements the DRMAA2 Go language binding based on top of
// Univa's DRMAA2 C API implementation. Should work also on other implementations
// which are compatible to the DRMAA2 standard.
// Please consult the DRMAA2 standard documents for more detailed
// information (http://www.ogf.org). More examples will also be
// published on my blog at http://www.gridengine.eu.
package drmaa2

import (
	"fmt"
	"log"
	"time"
	"unsafe"
)

/*
 #cgo LDFLAGS: -ldrmaa2 -O2 -g
 #include <stdio.h>
 #include <stdlib.h>
 #include <stddef.h>
 #include "drmaa2.h"

drmaa2_j malloc_job() {
   drmaa2_j job = (drmaa2_j) malloc(sizeof(drmaa2_j_s));
   job->id = NULL;
   job->session_name = NULL;
   return job;
}

drmaa2_jarray malloc_array_job() {
   drmaa2_jarray ja = (drmaa2_jarray) malloc(sizeof(drmaa2_jarray_s));
   ja->id = NULL;
   ja->session_name = NULL;
   ja->job_list = NULL;
   return ja;
}

drmaa2_jtemplate malloc_jtemplate() {
   drmaa2_jtemplate jt = (drmaa2_jtemplate) malloc(sizeof(drmaa2_jtemplate_s));
   jt->remoteCommand = DRMAA2_UNSET_STRING;
   jt->args = DRMAA2_UNSET_LIST;
   jt->submitAsHold = DRMAA2_UNSET_BOOL;
   jt->rerunnable = DRMAA2_UNSET_BOOL;
   jt->jobEnvironment = DRMAA2_UNSET_DICT;
   jt->workingDirectory = DRMAA2_UNSET_STRING;
   jt->jobCategory = DRMAA2_UNSET_STRING;
   jt->email = DRMAA2_UNSET_LIST;
   jt->emailOnStarted = DRMAA2_UNSET_BOOL;
   jt->emailOnTerminated = DRMAA2_UNSET_BOOL;
   jt->jobName = DRMAA2_UNSET_STRING;
   jt->inputPath = DRMAA2_UNSET_STRING;
   jt->outputPath = DRMAA2_UNSET_STRING;
   jt->errorPath = DRMAA2_UNSET_STRING;
   jt->joinFiles = DRMAA2_UNSET_BOOL;
   jt->reservationId = DRMAA2_UNSET_STRING;
   jt->queueName = DRMAA2_UNSET_STRING;
   jt->minSlots = DRMAA2_UNSET_NUM;
   jt->maxSlots = DRMAA2_UNSET_NUM;
   jt->priority = DRMAA2_UNSET_NUM;
   jt->candidateMachines = DRMAA2_UNSET_LIST;
   jt->minPhysMemory = DRMAA2_UNSET_NUM;
   jt->machineOS = DRMAA2_UNSET_ENUM;
   jt->machineArch = DRMAA2_UNSET_ENUM;
   jt->startTime = DRMAA2_UNSET_TIME;
   jt->deadlineTime = DRMAA2_UNSET_TIME;
   jt->stageInFiles = DRMAA2_UNSET_DICT;
   jt->stageOutFiles = DRMAA2_UNSET_DICT;
   jt->resourceLimits = DRMAA2_UNSET_DICT;
   jt->accountingId = DRMAA2_UNSET_STRING;
   jt->implementationSpecific = DRMAA2_UNSET_STRING;
   return jt;
}
*/
import "C"

// Interface definitions

// StructType is a type which represents the type of
// an extensible structure.
type StructType int

const (
	JobTemplateType = iota
	JobInfoType
	ReservationTemplateType
	ReservationInfoType
	QueueInfoType
	MachineInfoType
	NotificationType
)

// Extension is a struct which is embedded in DRMAA2 objects
// which are extensible.
type Extension struct {
	SType         StructType        // Stores the type of the struct
	Internal      unsafe.Pointer    // Enhancmement of C struct
	ExtensionList map[string]string // stores the extension requests as string
}

// Drmaa2Extensible is an interface which defines functions used to
// interact with extensible data structures (JobTemplate, JobInfo etc.).
type Drmaa2Extensible interface {
	// Lists all implementation specific key names for
	// a particular DRMAA2 extensible data type
	ListExtensions() []string
	DescribeExtension(string) string
	SetExtension(string) error
	GetExtension() string
	// points to data structure extension from C struct
}

// listExtension calls the C function for listing implementation specific
// enhancements for an object defined by the argument of the
// function.
func listExtensions(t StructType) []string {
	var clist C.drmaa2_string_list
	switch t {
	case JobTemplateType:
		clist = C.drmaa2_jtemplate_impl_spec()
	case JobInfoType:
		clist = C.drmaa2_jinfo_impl_spec()
	case ReservationTemplateType:
		clist = C.drmaa2_rtemplate_impl_spec()
	case QueueInfoType:
		clist = C.drmaa2_queueinfo_impl_spec()
	case MachineInfoType:
		clist = C.drmaa2_machineinfo_impl_spec()
	case ReservationInfoType:
		clist = C.drmaa2_rinfo_impl_spec()
	case NotificationType:
		clist = C.drmaa2_notification_impl_spec()
	default:
	}
	// cast string list in generic list type
	// since this is expected by the free function
	clistp := C.drmaa2_list(clist)
	defer C.drmaa2_list_free(&clistp)
	// Create a Go slice out of the string list.
	extensions := convertCStringListToGo(clist)
	return extensions
}

// ListExtensions returns a string list containing all implementation specific
// extensions of the JobTemplate object.
func (jt *JobTemplate) ListExtensions() []string {
	return listExtensions(JobTemplateType)
}

// ListExtensions returns a string list containing all implementation specific
// extensions of the Machine object.
func (m *Machine) ListExtensions() []string {
	return listExtensions(MachineInfoType)
}

// ListExtensions returns a string list containing all implementation specific
// extensions of the Queue object.
func (q *Queue) ListExtensions() []string {
	return listExtensions(QueueInfoType)
}

// ListExtensions returns a string list containing all implementation specific
// extensions of the JobInfo object.
func (ji *JobInfo) ListExtensions() []string {
	return listExtensions(JobInfoType)
}

func (ext *Extension) describeExtension(t StructType, extensionName string) (string, error) {
	if ext.Internal != nil {
		cdesc := C.drmaa2_describe_attribute(ext.Internal,
			C.CString(extensionName))
		if cdesc != nil {
			defer C.drmaa2_string_free(&cdesc)
			return C.GoString(cdesc), nil
		}
		return "", makeLastError()
	}
	// pointer to extension data structure is not set,
	// therefore it is allocated then used for the C
	// call and then thrown away - don't like it
	var description C.drmaa2_string

	switch t {
	case JobInfoType:
		jt := C.drmaa2_jtemplate_create()
		description = C.drmaa2_describe_attribute(jt.implementationSpecific,
			C.CString(extensionName))
		C.drmaa2_jtemplate_free(&jt)
	// TODO -> other types
	default:
		fmt.Println("Unimplemented")
	}

	if description != nil {
		defer C.drmaa2_string_free(&description)
		return C.GoString(description), nil
	}

	return "", makeLastError()
}

// DescribeExtension returns the description of an implementation specific
// JobTemplate extension as a string.
func (jt *JobTemplate) DescribeExtension(extensionName string) (string, error) {
	// good candidate for an init function in the session manager
	return jt.describeExtension(JobTemplateType, extensionName)
}

// TODO MachineInfo / Queue / JobInfo etc.

// checks if a certain extension exists for a given type
func extensionExists(t StructType, ext string) bool {
	// TODO expensive - better store available extensions
	// here a DRMAA2 init could be really useful
	extensions := listExtensions(t)
	for _, e := range extensions {
		if e == ext {
			return true
		}
	}
	return false
}

// Sets a DRM specific extension to a value
func (ext *Extension) setExtension(t StructType, extension, value string) error {
	if extensionExists(t, extension) {
		if ext.ExtensionList == nil {
			ext.ExtensionList = make(map[string]string)
		}
		ext.ExtensionList[extension] = value
		return nil
	}
	return makeError("Extension not supported", UnsupportedAttribute)
}

// SetExtension adds an vendor specific attribute to the extensible structure.
func (jt *JobTemplate) SetExtension(extension, value string) error {
	return jt.setExtension(JobTemplateType, extension, value)
}

// SetExtension adds an vendor specific attribute to the extensible structure.
func (m *Machine) SetExtension(extension, value string) error {
	return m.setExtension(MachineInfoType, extension, value)
}

// SetExtension adds an vendor specific attribute to the extensible structure.
func (ji *JobInfo) SetExtension(extension, value string) error {
	return ji.setExtension(JobInfoType, extension, value)
}

// SetExtension adds an vendor specific attribute to the extensible structure.
func (q *Queue) SetExtension(extension, value string) error {
	return q.setExtension(QueueInfoType, extension, value)
}

// TODO the other extensions: notification / reservation info / template

// set the Go extension into the real object
// (for example when running the job)
func setExtensionsIntoCObject(ptr unsafe.Pointer, elist map[string]string) {
	for key, value := range elist {
		C.drmaa2_set_instance_value(ptr, C.CString(key), C.CString(value))
	}
}

// GetExtension returns an extension of a extensible struct by name.
func (e *Extension) GetExtension(extension string) (string, error) {
	// check if any extension is stored in the Go struct
	if e.ExtensionList != nil {
		if value, ok := e.ExtensionList[extension]; ok == true {
			return value, nil
		}
		return "", makeError("Extension not found", UnsupportedAttribute)
	}
	return "", makeError("Extension not found", UnsupportedAttribute)
}

// Version is a DRMAA2 version type which consists of a major version
// and a minor version.
type Version struct {
	Major string
	Minor string
}

func (v *Version) String() string {
	return fmt.Sprintf("%s.%s", v.Major, v.Minor)
}

// Special timeout value: Don't wait
const ZeroTime = int64(C.DRMAA2_ZERO_TIME)

// Special timeout value: Wait probably infinitly
const InfiniteTime = int64(C.DRMAA2_INFINITE_TIME)

// Special time value: Time or date not set
const UnsetTime = int64(C.DRMAA2_UNSET_TIME)

// Capability is a type which represents the availability of optional
// functionality of the DRMAA2 implementation. Option functionality is
// defined by the DRMAA2 standard but not mandatory to implement.
type Capability int

const (
	AdvanceReservation = iota
	ReserveSlots
	Callback
	BulkJobsMaxParallel
	JtEmail
	JtStaging
	JtDeadline
	JtMaxSlots
	JtAccountingId
	RtStartNow
	RtDuration
	RtMachineOS
	RtMachineArch
)

// maybe not needed
var capCMap = map[C.drmaa2_capability]Capability{
	C.DRMAA2_ADVANCE_RESERVATION:   AdvanceReservation,
	C.DRMAA2_RESERVE_SLOTS:         ReserveSlots,
	C.DRMAA2_CALLBACK:              Callback,
	C.DRMAA2_BULK_JOBS_MAXPARALLEL: BulkJobsMaxParallel,
	C.DRMAA2_JT_EMAIL:              JtEmail,
	C.DRMAA2_JT_STAGING:            JtStaging,
	C.DRMAA2_JT_DEADLINE:           JtDeadline,
	C.DRMAA2_JT_MAXSLOTS:           JtMaxSlots,
	C.DRMAA2_JT_ACCOUNTINGID:       JtAccountingId,
	C.DRMAA2_RT_STARTNOW:           RtStartNow,
	C.DRMAA2_RT_DURATION:           RtDuration,
	C.DRMAA2_RT_MACHINEOS:          RtMachineOS,
	C.DRMAA2_RT_MACHINEARCH:        RtMachineArch,
}

var capMap = map[Capability]C.drmaa2_capability{
	AdvanceReservation:  C.DRMAA2_ADVANCE_RESERVATION,
	ReserveSlots:        C.DRMAA2_RESERVE_SLOTS,
	Callback:            C.DRMAA2_CALLBACK,
	BulkJobsMaxParallel: C.DRMAA2_BULK_JOBS_MAXPARALLEL,
	JtEmail:             C.DRMAA2_JT_EMAIL,
	JtStaging:           C.DRMAA2_JT_STAGING,
	JtDeadline:          C.DRMAA2_JT_DEADLINE,
	JtMaxSlots:          C.DRMAA2_JT_MAXSLOTS,
	JtAccountingId:      C.DRMAA2_JT_ACCOUNTINGID,
	RtStartNow:          C.DRMAA2_RT_STARTNOW,
	RtDuration:          C.DRMAA2_RT_DURATION,
	RtMachineOS:         C.DRMAA2_RT_MACHINEOS,
	RtMachineArch:       C.DRMAA2_RT_MACHINEARCH,
}

// DRMAA2 error ID
type ErrorId int

const (
	Success ErrorId = iota
	DeniedByDrms
	DrmCommunication
	TryLater
	SessionManagement
	Timeout
	Internal
	InvalidArgument
	InvalidSession
	InvalidState
	OutOfResource
	UnsupportedAttribute
	UnsupportedOperation
	ImplementationSpecific
	LastError
)

// Maps a C drmaa2_error type into a Go ErrorId
var errorIdMap = map[C.drmaa2_error]ErrorId{
	C.DRMAA2_SUCCESS:                 Success,
	C.DRMAA2_DENIED_BY_DRMS:          DeniedByDrms,
	C.DRMAA2_DRM_COMMUNICATION:       DrmCommunication,
	C.DRMAA2_TRY_LATER:               TryLater,
	C.DRMAA2_SESSION_MANAGEMENT:      SessionManagement,
	C.DRMAA2_TIMEOUT:                 Timeout,
	C.DRMAA2_INTERNAL:                Internal,
	C.DRMAA2_INVALID_ARGUMENT:        InvalidArgument,
	C.DRMAA2_INVALID_SESSION:         InvalidSession,
	C.DRMAA2_INVALID_STATE:           InvalidState,
	C.DRMAA2_OUT_OF_RESOURCE:         OutOfResource,
	C.DRMAA2_UNSUPPORTED_ATTRIBUTE:   UnsupportedAttribute,
	C.DRMAA2_UNSUPPORTED_OPERATION:   UnsupportedOperation,
	C.DRMAA2_IMPLEMENTATION_SPECIFIC: ImplementationSpecific,
	C.DRMAA2_LASTERROR:               LastError,
}

// CPU architecture types
type CPU int

const (
	OtherCPU CPU = iota
	Alpha
	ARM
	ARM64
	Cell
	PA_RISC
	PA_RISC64
	x86
	x64
	IA_64
	MIPS
	MIPS64
	PowerPC
	PowerPC64
	SPARC
	SPARC64
)

var cpuMap = map[C.drmaa2_cpu]CPU{
	C.DRMAA2_OTHER_CPU: OtherCPU,
	C.DRMAA2_ALPHA:     Alpha,
	C.DRMAA2_ARM:       ARM,
	C.DRMAA2_ARM64:     ARM64,
	C.DRMAA2_CELL:      Cell,
	C.DRMAA2_PARISC:    PA_RISC,
	C.DRMAA2_PARISC64:  PA_RISC64,
	C.DRMAA2_X86:       x86,
	C.DRMAA2_X64:       x64,
	C.DRMAA2_IA64:      IA_64,
	C.DRMAA2_MIPS:      MIPS,
	C.DRMAA2_MIPS64:    MIPS64,
	C.DRMAA2_PPC:       PowerPC,
	C.DRMAA2_PPC64:     PowerPC64,
	C.DRMAA2_SPARC:     SPARC,
	C.DRMAA2_SPARC64:   SPARC64,
}

func (cpu CPU) String() string {
	switch cpu {
	case OtherCPU:
		return "OtherCPU"
	case Alpha:
		return "Alpha"
	case ARM:
		return "ARM"
	case ARM64:
		return "ARM64"
	case Cell:
		return "Cell"
	case PA_RISC:
		return "PA_RISC"
	case PA_RISC64:
		return "PA_RISC64"
	case x86:
		return "x86"
	case x64:
		return "x64"
	case IA_64:
		return "IA_64"
	case MIPS:
		return "MIPS"
	case MIPS64:
		return "MIPS64"
	case PowerPC:
		return "PowerPC"
	case SPARC:
		return "SPARC"
	case SPARC64:
		return "SPARC64"
	}
	return "Unknown"
}

// OS is the operating system type.
type OS int

const (
	OtherOS OS = iota
	AIX
	BSD
	Linux
	HPUX
	IRIX
	MacOS
	SunOS
	TRU64
	UnixWare
	Win
	WinNT
)

// An OS struct needs to be printable.
func (os OS) String() string {
	switch os {
	case OtherOS:
		return "OtherOS"
	case AIX:
		return "AIX"
	case BSD:
		return "BSD"
	case Linux:
		return "Linux"
	case HPUX:
		return "HPUX"
	case IRIX:
		return "IRIX"
	case MacOS:
		return "MacOS"
	case SunOS:
		return "SunOS"
	case TRU64:
		return "TRU64"
	case UnixWare:
		return "UnixWare"
	case Win:
		return "Win"
	case WinNT:
		return "WinNT"
	}
	return "Unknown"
}

var osMap = map[C.drmaa2_os]OS{
	C.DRMAA2_OTHER_OS: OtherOS,
	C.DRMAA2_AIX:      AIX,
	C.DRMAA2_BSD:      BSD,
	C.DRMAA2_LINUX:    Linux,
	C.DRMAA2_HPUX:     HPUX,
	C.DRMAA2_IRIX:     IRIX,
	C.DRMAA2_MACOS:    MacOS,
	C.DRMAA2_SUNOS:    SunOS,
	C.DRMAA2_TRU64:    TRU64,
	C.DRMAA2_UNIXWARE: UnixWare,
	C.DRMAA2_WIN:      Win,
	C.DRMAA2_WINNT:    WinNT,
}

// Job States
type JobState int

const (
	Unset JobState = iota
	Undetermined
	Queued
	QueuedHeld
	Running
	Suspended
	Requeued
	RequeuedHeld
	Done
	Failed
)

// Implements the Stringer interface
func (js JobState) String() string {
	switch js {
	case Undetermined:
		return "Undetermined"
	case Queued:
		return "Queued"
	case QueuedHeld:
		return "QueuedHeld"
	case Running:
		return "Running"
	case Suspended:
		return "Suspended"
	case Requeued:
		return "Requeued"
	case RequeuedHeld:
		return "RequeuedHeld"
	case Done:
		return "Done"
	case Failed:
		return "Failed"
	}
	return "Unset"
}

// TODO UNSET_ENUM can't be mapped to a JobState
// because it is not from that type -> use
// Undetermined instead
var jobStateMap = map[C.drmaa2_jstate]JobState{
	C.DRMAA2_UNDETERMINED:  Undetermined,
	C.DRMAA2_QUEUED:        Queued,
	C.DRMAA2_QUEUED_HELD:   QueuedHeld,
	C.DRMAA2_RUNNING:       Running,
	C.DRMAA2_SUSPENDED:     Suspended,
	C.DRMAA2_REQUEUED_HELD: RequeuedHeld,
	C.DRMAA2_DONE:          Done,
	C.DRMAA2_FAILED:        Failed,
}

// convertGoStateToC returns the DRMAA2 C API state which
// is the equivalent to the Go API counterpart.
func convertGoStateToC(s JobState) (state C.drmaa2_jstate) {
	switch s {
	case Undetermined:
		return C.DRMAA2_UNDETERMINED
	case Queued:
		return C.DRMAA2_QUEUED
	case QueuedHeld:
		return C.DRMAA2_QUEUED_HELD
	case Running:
		return C.DRMAA2_RUNNING
	case Suspended:
		return C.DRMAA2_SUSPENDED
	case RequeuedHeld:
		return C.DRMAA2_REQUEUED_HELD
	case Done:
		return C.DRMAA2_DONE
	case Failed:
		return C.DRMAA2_FAILED
	}
	return C.DRMAA2_UNDETERMINED
}

// DRMAA2 error (implements GO Error interface).
type Error struct {
	Message string
	Id      ErrorId
}

// The DRMAA2 Error implements the Error interface.
func (ce Error) Error() string {
	return ce.Message
}

// Implement the Stringer interface for an drmaa2.Error
func (ce Error) String() string {
	return ce.Message
}

// Intenal function which creats an GO DRMAA2 error.
func makeError(msg string, id ErrorId) Error {
	var ce Error
	ce.Message = msg
	ce.Id = id
	return ce
}

func makeLastError() *Error {
	cerr := C.drmaa2_lasterror_text()
	defer C.free(unsafe.Pointer(cerr))
	msg := C.GoString(cerr)
	id := C.drmaa2_lasterror()
	err := makeError(msg, errorIdMap[id])
	return &err
}

// TODO(dg) A Create Method which initializes the values and
// also does initialization about capabilities,
// versions etc. ?!?
type SessionManager struct {
	//drmsName     string
	//drmsVersion  string // type Version
	//drmaaName    string
	//drmaaVersion string // type Version
}

// MonitoringSession is a struct which represents a DRMAA2
// monitoring session (for cluster monitoring).
type MonitoringSession struct {
	name string            // internal
	ms   C.drmaa2_msession // pointer to C drmaa2 session type
}

// JobSession is a struct which represents a DRMAA2 job session
// for job submission and job management.
type JobSession struct {
	Name string            `json:"name"` // public name of job session
	js   C.drmaa2_jsession // pointer to C drmaa2 job session type
}

// ReservationSession is a struct which represents a DRMAA reservation
// session for handling advance reservations. The support for reservation
// sessions is optional in the DRMAA2 standard.
type ReservationSession struct {
	Name string `json:"name"`
	rs   C.drmaa2_rsession
}

// ReservationInfo is a struct which represents a reservation.
type ReservationInfo struct {
	ReservationId        string    `json:"reservationId"`
	ReservationName      string    `json:"reservationName"`
	ReservationStartTime time.Time `json:"reservationStartTime"`
	ReservationEndTime   time.Time `json:"reservationEndTime"`
	ACL                  []string  `json:"acl"`
	ReservedSlots        int64     `json:"reservedSlots"`
	ReservedMachines     []string  `json:"reservedMachines"`
}

// Job is a struct which represents a job in the DRM system (part of a
// JobSession or a Monitoring session).
type Job struct {
	// job is private implementation specific (see struct drmaa2_j_s)
	id           string
	session_name string
}

// JobInfor is a struct which represents the current state of a job.
type JobInfo struct {
	// reference to the void* pointer which
	// is used for extensions
	Extension         `xml:"-" json:"-"`
	Id                string        `json:"id"`
	ExitStatus        int           `json:"exitStatus"`
	TerminatingSignal string        `json:"terminationSignal"`
	Annotation        string        `json:"annotation"`
	State             JobState      `json:"state"`
	SubState          string        `json:"subState"`
	AllocatedMachines []string      `json:"allocatedMachines"`
	SubmissionMachine string        `json:"submissionMachine"`
	JobOwner          string        `json:"jobOwner"`
	Slots             int64         `json:"slots"`
	QueueName         string        `json:"queueName"`
	WallclockTime     time.Duration `json:"wallockTime"`
	CPUTime           int64         `json:"cpuTime"`
	SubmissionTime    time.Time     `json:"submissionTime"`
	DispatchTime      time.Time     `json:"dispatchTime"`
	FinishTime        time.Time     `json:"finishTime"`
}

// CreateJobInfo creates a JobInfo object where all values are initialized
// with UNSET (needed in order to differentiate if a value is
// not set or 0).
func CreateJobInfo() (ji JobInfo) {
	// strings are unset with ""
	ji.ExitStatus = C.DRMAA2_UNSET_NUM
	// slices are unset with nil
	ji.Slots = C.DRMAA2_UNSET_NUM
	// WallclockTime is unset with 0
	ji.CPUTime = C.DRMAA2_UNSET_TIME
	ji.State = Unset
	// TODO Unset for Go Time type...
	return ji
}

type ArrayJob struct {
	// needed for suspend / resume ...
	aj          C.drmaa2_jarray
	id          string
	jobs        []Job
	sessionName string
	jt          JobTemplate
}

type Queue struct {
	Extension `xml:"-" json:"-"`
	Name      string `xml:"name"`
}

type Machine struct {
	Extension      `xml:"-" json:"-"`
	Name           string  `json:"name"`
	Available      bool    `json:"available"`
	Sockets        int64   `json:"sockets"`
	CoresPerSocket int64   `json:"coresPerSocket"`
	ThreadsPerCore int64   `json:"threadsPerCore"`
	Load           float64 `json:"load"`
	PhysicalMemory int64   `json:"physicalMemory"`
	VirtualMemory  int64   `json:"virtualMemory"`
	Architecture   CPU     `json:"architecture"`
	OSVersion      Version `json:"osVersion"`
	OS             OS      `json:"os"`
}

type JobTemplate struct {
	Extension         `xml:"-" json:"-"`
	RemoteCommand     string            `json:"remoteCommand"`
	Args              []string          `json:"args"`
	SubmitAsHold      bool              `json:"submitAsHold"`
	ReRunnable        bool              `json:"reRunnable"`
	JobEnvironment    map[string]string `json:"jobEnvironment"`
	WorkingDirectory  string            `json:"workingDirectory"`
	JobCategory       string            `json:"jobCategory"`
	Email             []string          `json:"email"`
	EmailOnStarted    bool              `json:"emailOnStarted"`
	EmailOnTerminated bool              `json:"emailOnTerminated"`
	JobName           string            `json:"jobName"`
	InputPath         string            `json:"inputPath"`
	OutputPath        string            `json:"outputPath"`
	ErrorPath         string            `json:"errorPath"`
	JoinFiles         bool              `json:"joinFiles"`
	ReservationId     string            `json:"reservationId"`
	QueueName         string            `json:"queueName"`
	MinSlots          int64             `json:"minSlots"`
	MaxSlots          int64             `json:"maxSlots"`
	Priority          int64             `json:"priority"`
	CandidateMachines []string          `json:"candidateMachines"`
	MinPhysMemory     int64             `json:"minPhysMemory"`
	MachineOs         string            `json:"machineOs"`
	MachineArch       string            `json:"machineArch"`
	StartTime         time.Time         `json:"startTime"`
	DeadlineTime      time.Time         `json:"deadlineTime"`
	StageInFiles      map[string]string `json:"stageInFiles"`
	StageOutFiles     map[string]string `json:"stageOutFiles"`
	ResourceLimits    map[string]string `json:"resourceLimits"`
	AccountingId      string            `json:"accountingString"`
}

type ReservationTemplate struct {
	Extension         `xml:"-" json:"-"`
	Name              string        `json:"name"`
	StartTime         time.Time     `json:"startTime"`
	EndTime           time.Time     `json:"endTime"`
	Duration          time.Duration `json:"duration"`
	MinSlots          int64         `json:"minSlots"`
	MaxSlots          int64         `json:"maxSlots"`
	JobCategory       string        `json:"jobCategory"`
	UsersACL          []string      `json:"userACL"`
	CandidateMachines []string      `json:"candidateMachines"`
	MinPhysMemory     int64         `json:"minPhysMemory"`
	MachineOs         string        `json:"machineOs"`
	MachineArch       string        `json:"machineArch"`
}

type Reservation struct {
	SessionName   string              `json:"sessionName"`
	Contact       string              `json:"contact"`
	Template      ReservationTemplate `json:"template"`
	ReservationId string              `json:"reservationId"`
}

// this is needed since there is a difference between "" and nil
func convertGoStringToC(s string) C.drmaa2_string {
	if s != "" {
		return C.CString(s)
	}
	return nil
}

// Converts a JobTemplate in the C DRMAA2 equivalent
// and sets the values.
func convertGoJtemplateToC(jt JobTemplate) C.drmaa2_jtemplate {
	cjt := C.malloc_jtemplate()
	cjt.remoteCommand = convertGoStringToC(jt.RemoteCommand)
	cjt.args = C.drmaa2_string_list(convertGoListToC(jt.Args))
	cjt.submitAsHold = convertGoBoolToC(jt.SubmitAsHold)
	cjt.rerunnable = convertGoBoolToC(jt.ReRunnable)
	cjt.jobEnvironment = convertGoDictToC(jt.JobEnvironment)
	cjt.workingDirectory = convertGoStringToC(jt.WorkingDirectory)
	cjt.jobCategory = convertGoStringToC(jt.JobCategory)
	cjt.email = C.drmaa2_string_list(convertGoListToC(jt.Email))
	cjt.emailOnStarted = convertGoBoolToC(jt.EmailOnStarted)
	cjt.emailOnTerminated = convertGoBoolToC(jt.EmailOnTerminated)
	cjt.jobName = convertGoStringToC(jt.JobName)
	cjt.inputPath = convertGoStringToC(jt.InputPath)
	cjt.outputPath = convertGoStringToC(jt.OutputPath)
	cjt.errorPath = convertGoStringToC(jt.ErrorPath)
	cjt.joinFiles = convertGoBoolToC(jt.JoinFiles)
	cjt.reservationId = convertGoStringToC(jt.ReservationId)
	cjt.queueName = convertGoStringToC(jt.QueueName)
	// TODO initialize JobTemplate with UNSET values!
	if jt.MinSlots > 0 {
		cjt.minSlots = C.longlong(jt.MinSlots)
	}
	if jt.MaxSlots > 0 {
		cjt.maxSlots = C.longlong(jt.MaxSlots)
	}
	if jt.Priority != 0 {
		cjt.priority = C.longlong(jt.Priority)
	}
	cjt.candidateMachines = C.drmaa2_string_list(convertGoListToC(jt.CandidateMachines))
	if jt.MinPhysMemory > 0 {
		cjt.minPhysMemory = C.longlong(jt.MinPhysMemory)
	}
	// machineOs
	// machineArch
	// startTime
	// deadlineTime
	cjt.stageInFiles = convertGoDictToC(jt.StageInFiles)
	cjt.stageOutFiles = convertGoDictToC(jt.StageOutFiles)
	cjt.resourceLimits = convertGoDictToC(jt.ResourceLimits)
	cjt.accountingId = convertGoStringToC(jt.AccountingId)

	return cjt
}

// Converts a JobInfo in the C counterpart.
// Needs to be freed! One point is that values in ji
// need to be UNSET...
func convertGoJobInfoToC(ji JobInfo) C.drmaa2_jinfo {
	cji := C.drmaa2_jinfo_create()
	// TODO JobName is missing in JobInfo (DRMAA2 issue)
	cji.jobId = convertGoStringToC(ji.Id)
	if ji.ExitStatus != C.DRMAA2_UNSET_NUM {
		cji.exitStatus = C.int(ji.ExitStatus)
	}
	cji.terminatingSignal = convertGoStringToC(ji.TerminatingSignal)
	cji.annotation = convertGoStringToC(ji.Annotation)
	// TODO check spec
	if ji.State != Unset {
		cji.jobState = convertGoStateToC(ji.State)
	}
	cji.jobSubState = convertGoStringToC(ji.SubState)
	//cji.allocatedMachines = C.drmaa2_string_list(convertGoListToC(ji.AllocatedMachines))
	cji.submissionMachine = convertGoStringToC(ji.SubmissionMachine)
	cji.jobOwner = convertGoStringToC(ji.JobOwner)
	//cji.slots = C.longlong(ji.Slots)
	cji.queueName = convertGoStringToC(ji.QueueName)

	// TODO
	// cji.wallclockTime
	// cji.cpuTime
	// cji.submissionTime
	// cji.dispatchTime
	// cji.finishTime

	return cji
}

// Converts a element from a DRMAA2 list into
// the C counterpart and treat it like a void*
// pointer.
func convertListElement(element interface{}) unsafe.Pointer {
	switch element.(type) {
	case Job:
		return unsafe.Pointer(convertGoJobToC(element.(Job)))
	case string:
		return unsafe.Pointer(C.CString(element.(string)))
	default:
		// unexpected type
		log.Fatal("convertListElement unknown type")
	}
	return nil
}

// Data Type conversion
func convertCStringListToGo(cl C.drmaa2_string_list) []string {
	// TODO Cgecj if it is NULL
	length := int64(C.drmaa2_list_size(C.drmaa2_list(cl)))

	list := make([]string, length, length)
	for i := int64(0); i < length; i++ {
		element := C.GoString(C.drmaa2_string(C.drmaa2_list_get(C.drmaa2_list(cl), C.long(i))))
		list[i] = element
	}
	return list
}

// conertGoListToC converts a Go list into the C DRMAA2 counterpart
// which needs to be freed by the caller
func convertGoListToC(list interface{}) C.drmaa2_list {
	var l C.drmaa2_list
	switch list.(type) {
	case []Job:
		tlist := []Job(list.([]Job))
		l = C.drmaa2_list_create(C.DRMAA2_JOBLIST, nil)
		for _, e := range tlist {
			C.drmaa2_list_add(l, unsafe.Pointer(convertGoJobToC(e)))
		}
	case []string:
		tlist := []string(list.([]string))
		l = C.drmaa2_list_create(C.DRMAA2_STRINGLIST, nil)
		for _, e := range tlist {
			C.drmaa2_list_add(l, unsafe.Pointer(C.CString(e)))
		}
	default:
		// unexpected type
		log.Fatal("convertGoListToC: unexpected type")
	}
	return l
}

func convertGoBoolToC(value bool) C.drmaa2_bool {
	if value == true {
		return C.DRMAA2_TRUE
	}
	return C.DRMAA2_FALSE
}

func convertGoDictToC(dict map[string]string) C.drmaa2_dict {
	// just initialize it with NULL
	if dict == nil || len(dict) <= 0 {
		return nil
	}
	cdict := C.drmaa2_dict_create(nil)
	for k, v := range dict {
		C.drmaa2_dict_set(C.drmaa2_dict(cdict), C.CString(k), C.CString(v))
	}
	return cdict
}

/* Helper for array jobs. */
func convertGoArrayJobToC(ja ArrayJob) C.drmaa2_jarray {
	caj := C.malloc_array_job()
	caj.id = C.CString(ja.id)
	caj.session_name = C.CString(ja.sessionName)
	return caj
}

func convertCArrayJobToGo(ja C.drmaa2_jarray) ArrayJob {
	var aj ArrayJob
	//aj.aj = ja
	aj.id = C.GoString(ja.id)
	aj.sessionName = C.GoString(ja.session_name)
	aj.jobs = convertCJobListToGo(ja.job_list)
	// add array job
	jt := C.drmaa2_jarray_get_job_template(ja)
	aj.jt = convertCJtemplateToGo(jt)
	return aj
}

/* Methods working on job. */
func convertCJobToGo(cj C.drmaa2_j) Job {
	var job Job
	job.id = C.GoString(cj.id)
	job.session_name = C.GoString(cj.session_name)
	return job
}

func convertGoJobToC(job Job) C.drmaa2_j {
	cjob := C.malloc_job()
	cjob.id = C.CString(job.id)
	cjob.session_name = C.CString(job.session_name)
	return cjob
}

func (job *Job) GetId() string {
	return job.id
}

// A job session name is a (per user) unique identifier
// of the job session. It is stored in the cluster scheduler
// or in the underlying DRMAA2 C implementation. It is
// persistent until it gets reaped by DestroyJobSession()
// method.
func (job *Job) GetSessionName() string {
	return job.session_name
}

func goBool(v C.drmaa2_bool) bool {
	if v == C.DRMAA2_TRUE {
		return true
	}
	return false
}

func goStringList(string_list C.drmaa2_string_list) []string {
	strings := make([]string, 0)
	if string_list != nil {
		size := (int64)(C.drmaa2_list_size((C.drmaa2_list)(string_list)))
		for i := (int64)(0); i < size; i++ {
			cstr := (*C.char)(C.drmaa2_list_get((C.drmaa2_list)(string_list), C.long(i)))
			strings = append(strings, C.GoString(cstr))
		}
	}
	return strings
}

func goOS(os C.drmaa2_os) OS {
	return osMap[os]
}

func goVersion(version C.drmaa2_version) (v Version) {
	if version == nil {
		v.Major = "0"
		v.Minor = "0"
		return v
	}
	v.Major = C.GoString(version.major)
	v.Minor = C.GoString(version.minor)
	return v
}

func goArchitecture(cpu C.drmaa2_cpu) CPU {
	return cpuMap[cpu]
}

func goJobState(state C.drmaa2_jstate) JobState {
	return jobStateMap[state]
}

// goTime reates a point in Time out of a C time stamp
func goTime(sec C.time_t) time.Time {
	// if time C.DRMAA2_UNSET_TIME
	return time.Unix((int64)(sec), (int64)(0))
}

// goDuration creates a Duration out of a C time in seconds
func goDuration(sec C.time_t) time.Duration {
	timeInSeconds := fmt.Sprintf("%ds", (int64)(sec))
	duration, _ := time.ParseDuration(timeInSeconds)
	return duration
}

// goJobInfo converts a C Job info in a Go Job Info object
func goJobInfo(cji C.drmaa2_jinfo) JobInfo {
	var jinfo JobInfo
	/* convert C job info into Go job info */
	ji := (C.drmaa2_jinfo_s)(*cji)
	jinfo.AllocatedMachines = goStringList(ji.allocatedMachines)
	jinfo.Annotation = C.GoString(ji.annotation)
	jinfo.CPUTime = (int64)(ji.cpuTime)
	jinfo.ExitStatus = (int)(ji.exitStatus)
	jinfo.Id = C.GoString(ji.jobId)
	jinfo.JobOwner = C.GoString(ji.jobOwner)
	jinfo.QueueName = C.GoString(ji.queueName)
	jinfo.Slots = (int64)(ji.slots)
	jinfo.State = goJobState(ji.jobState)
	jinfo.SubState = C.GoString(ji.jobSubState)
	jinfo.SubmissionTime = goTime(ji.submissionTime)
	jinfo.SubmissionMachine = C.GoString(ji.submissionMachine)
	jinfo.TerminatingSignal = C.GoString(ji.terminatingSignal)
	jinfo.WallclockTime = goDuration(ji.wallclockTime)
	jinfo.DispatchTime = goTime(ji.dispatchTime)
	jinfo.FinishTime = goTime(ji.finishTime)
	return jinfo
}

// helper function for converting c jtemplate to go
func convertCJtemplateToGo(t C.drmaa2_jtemplate) JobTemplate {
	var jt JobTemplate
	jt.AccountingId = C.GoString(t.accountingId)
	jt.Args = goStringList(t.args)
	jt.EmailOnStarted = goBool(t.emailOnStarted)
	jt.EmailOnTerminated = goBool(t.emailOnTerminated)
	// TODO dict
	jt.ErrorPath = C.GoString(t.errorPath)
	jt.InputPath = C.GoString(t.inputPath)
	jt.JobCategory = C.GoString(t.jobCategory)
	// TOOD jt.JobEnvironment
	jt.JobName = C.GoString(t.jobName)
	jt.JoinFiles = goBool(t.joinFiles)
	//jt.MachineArch = C.GoString(t.machineArch)
	jt.MaxSlots = (int64)(t.maxSlots)
	jt.MinPhysMemory = (int64)(t.minPhysMemory)
	jt.MinSlots = (int64)(t.minSlots)
	jt.OutputPath = C.GoString(t.outputPath)
	jt.Priority = (int64)(t.priority)
	jt.QueueName = C.GoString(t.queueName)
	jt.ReRunnable = goBool(t.rerunnable)
	jt.RemoteCommand = C.GoString(t.remoteCommand)
	jt.ReservationId = C.GoString(t.reservationId)
	// jt.ResourceLimits
	// jt.StageInFiles
	// jt.StageOutFiles
	jt.SubmitAsHold = goBool(t.submitAsHold)
	jt.WorkingDirectory = C.GoString(t.workingDirectory)
	// jt.machineOs convert ennum
	return jt
}

// Returns the JobTemplate used to submit the job.
func (job *Job) GetJobTemplate() (*JobTemplate, error) {
	cjob := convertGoJobToC(*job)

	if cjob == nil {
		return nil, makeLastError()
	}
	defer C.drmaa2_j_free(&cjob)

	cjt := C.drmaa2_j_get_jt(cjob)
	// TODO convert C job template into Go jobtemplate
	if cjt != nil {
		defer C.drmaa2_jtemplate_free(&cjt)
		jt := convertCJtemplateToGo(cjt)
		return &jt, nil
	}
	return nil, makeLastError()
}

// GetState returns the current state of the job.
func (job *Job) GetState() JobState {
	if ji, err := job.GetJobInfo(); err == nil {
		return ji.State
	}
	return Undetermined
}

// Creates a JobInfo object from the job containing
// more detailed information about the job.
func (job *Job) GetJobInfo() (*JobInfo, error) {
	cjob := convertGoJobToC(*job)
	if cjob == nil {
		return nil, makeLastError()
	}
	defer C.drmaa2_j_free(&cjob)

	cji := C.drmaa2_j_get_info(cjob)
	if cji == nil {
		return nil, makeLastError()
	}
	defer C.drmaa2_jinfo_free(&cji)

	/* convert C job info into Go job info */
	jinfo := goJobInfo(cji)

	return &jinfo, nil
}

// internal operations on job
type modop int

const (
	suspend = iota
	resume
	hold
	release
	terminate
)

func (job *Job) modify(operation modop) error {
	cjob := convertGoJobToC(*job)
	var ret C.drmaa2_error

	switch operation {
	case suspend:
		ret = C.drmaa2_j_suspend(cjob)
	case resume:
		ret = C.drmaa2_j_resume(cjob)
	case hold:
		ret = C.drmaa2_j_hold(cjob)
	case release:
		ret = C.drmaa2_j_release(cjob)
	case terminate:
		ret = C.drmaa2_j_terminate(cjob)
	}
	defer C.drmaa2_j_free(&cjob)
	if ret != C.DRMAA2_SUCCESS {
		return makeLastError()
	}
	return nil
}

// Stops a job / process from beeing executed.
func (job *Job) Suspend() error {
	return job.modify(suspend)
}

// Resumes a job / process to be runnable again.
func (job *Job) Resume() error {
	return job.modify(resume)
}

// Puts the job into an hold state so that it is not
// scheduled. If the job is already running it continues
// to run and the hold state becomes only effectice when
// the job is rescheduled.
func (job *Job) Hold() error {
	return job.modify(hold)
}

// Releases the job from hold state so that it will
// be schedulable.
func (job *Job) Release() error {
	return job.modify(release)
}

// Terminate tells the resource manager to kill the job.
func (job *Job) Terminate() error {
	return job.modify(terminate)
}

// Blocking wait until the job is started. The timeout
// prefents that the call is blocking endlessly. Special
// timeouts are available by the constants InfiniteTime
// and ZeroTime.
func (job *Job) WaitStarted(timeout int64) error {
	cjob := convertGoJobToC(*job)
	//defer C.drmaa2_j_free(&cjob)
	err := C.drmaa2_j_wait_started(cjob, (C.time_t)(timeout))
	if err != C.DRMAA2_SUCCESS {
		return makeLastError()
	}
	return nil
}

// Waits until the job goes into one of the finished states.
// The timeout specifies the maximum time to wait.
// If no timeout is required use the constant InfiniteTime.
func (job *Job) WaitTerminated(timeout int64) error {
	cjob := convertGoJobToC(*job)
	defer C.drmaa2_j_free(&cjob)
	err := C.drmaa2_j_wait_terminated(cjob, (C.time_t)(timeout))
	if err != C.DRMAA2_SUCCESS {
		return makeLastError()
	}
	return nil
}

// CreateJobSession creates a new persistent job session and opens it. The
// returned JobSession object contains a reference to a DRMAA2 C jobsession
// object and hence needs to be freed manually.
func (sm *SessionManager) CreateJobSession(sessionName, contact string) (*JobSession, error) {
	var js JobSession
	// convert parameters
	name := C.CString(sessionName)
	defer C.free(unsafe.Pointer(name))
	if contact == "" {
		// default in Univa Grid Engine case
		ctct := C.drmaa2_string(nil)
		js.js = C.drmaa2_create_jsession(name, ctct)
	} else {
		ctct := C.drmaa2_string(C.CString(contact))
		js.js = C.drmaa2_create_jsession(name, ctct)
	}
	// convert error back to Go
	if js.js == nil {
		// an error happended - create an error
		return nil, makeLastError()
	}
	// job session needs to be freed from caller
	return &js, nil
}

// Creates a ReservationSession by name and contact string.
func (sm *SessionManager) CreateReservationSession(sessionName, contact string) (rs *ReservationSession, err error) {
	return rs, nil
}

// Opens a MonitoringSession by name. Usually the name is ignored.
func (sm *SessionManager) OpenMonitoringSession(sessionName string) (*MonitoringSession, error) {
	var ms MonitoringSession
	if sessionName != "" {
		snp := C.CString(sessionName)
		defer C.free(unsafe.Pointer(snp))
		// how to convert msession to a Go representation?
		C.drmaa2_open_msession(snp)
	} else {
		ms.ms = C.drmaa2_open_msession(nil)
		return &ms, nil
	}
	if ms.ms == nil {
		// an error happend -> get error and return as string
		return nil, makeLastError()
	}
	ms.name = sessionName
	return &ms, nil
}

// Closes the MonitoringSession. TODO as method or as part of SessionManager?
func (ms *MonitoringSession) CloseMonitoringSession() error {
	err_cstr := C.drmaa2_close_msession(ms.ms)
	if err_cstr == C.DRMAA2_SUCCESS {
		C.drmaa2_msession_free(&ms.ms)
		return nil
	}
	C.drmaa2_msession_free(&ms.ms)
	return makeLastError()
}

func convertCJobListToGo(jlist C.drmaa2_j_list) []Job {
	if jlist == nil {
		return nil
	}
	jl := (C.drmaa2_list)(jlist)
	count := (int64)(C.drmaa2_list_size(jl))
	// ...
	jobs := make([]Job, 0)
	for i := (int64)(0); i < count; i++ {
		cjob := (C.drmaa2_j)(C.drmaa2_list_get(jl, C.long(i)))
		if cjob == nil {
			continue
		}
		// copy C implementation specific
		// job struct values -> therefore we need
		// access to Grid Engine internal header file
		var j Job
		cj := (C.drmaa2_j_s)(*cjob)
		j.id = C.GoString(cj.id)
		j.session_name = C.GoString(cj.session_name)
		jobs = append(jobs, j)
	}
	return jobs
}

// Creates a slice of Queues based on C queue list.
func createQueueList(ql C.drmaa2_list) []Queue {
	if ql == nil {
		return nil
	}
	queues := make([]Queue, 0)
	count := (int64)(C.drmaa2_list_size(ql))
	// ...
	for i := (int64)(0); i < count; i++ {
		cq := (C.drmaa2_queueinfo)(C.drmaa2_list_get(ql, C.long(i)))
		if cq == nil {
			continue
		}
		// copy public visible string name
		var q Queue
		cqi := *cq
		q.Name = C.GoString(cqi.name)
		queues = append(queues, q)
	}
	return queues
}

func createMachineList(ml C.drmaa2_list) []Machine {
	if ml == nil {
		return nil
	}
	machines := make([]Machine, 0)
	count := (int64)(C.drmaa2_list_size(ml))
	// ...
	for i := (int64)(0); i < count; i++ {
		mi := (C.drmaa2_machineinfo)(C.drmaa2_list_get(ml, C.long(i)))
		if mi == nil {
			continue
		}
		// copy public visible string name
		var m Machine
		cmi := *mi
		m.Name = C.GoString(cmi.name)
		m.Available = goBool(cmi.available)
		m.Architecture = goArchitecture(cmi.machineArch)
		m.Sockets = (int64)(cmi.sockets)
		m.CoresPerSocket = (int64)(cmi.coresPerSocket)
		m.ThreadsPerCore = (int64)(cmi.threadsPerCore)
		m.PhysicalMemory = (int64)(cmi.physMemory)
		m.VirtualMemory = (int64)(cmi.virtMemory)
		m.OS = goOS(cmi.machineOS)
		m.Load = (float64)(cmi.load)
		m.OSVersion = goVersion(cmi.machineOSVersion)
		machines = append(machines, m)
	}
	return machines
}

// GetAllJobs returns a slice of jobs currently visible in the monitoring session.
// The JobInfo parameter specifies a filter for the job. For instance
// when a certain job number is set in the JobInfo object, then
func (ms *MonitoringSession) GetAllJobs(ji *JobInfo) (jobs []Job, err error) {
	// Create the job filter
	var cji C.drmaa2_jinfo
	if ji != nil {
		cji = convertGoJobInfoToC(*ji)
		defer C.drmaa2_jinfo_free(&cji)
	} else {
		cji = nil
	}
	cjlist := (C.drmaa2_j_list)(C.drmaa2_msession_get_all_jobs(ms.ms, cji))
	if cjlist == nil {
		return nil, makeLastError()
	}
	jl := convertCJobListToGo(cjlist)
	jlist := (C.drmaa2_list)(cjlist)
	C.drmaa2_list_free(&jlist)
	return jl, nil
}

// GetlAllQueues returns all queues configured in the cluster in case the argument is
// nil. Otherwise as subset of the queues which matches the given names
// is returned.
func (ms *MonitoringSession) GetAllQueues(names []string) (queues []Queue, err error) {
	var arg C.drmaa2_string_list
	if names == nil {
		arg = nil
	} else {
		arg = C.drmaa2_string_list(convertGoListToC(names))
	}

	cqlist := (C.drmaa2_list)(C.drmaa2_msession_get_all_queues(ms.ms, arg))
	if cqlist == nil {
		return nil, makeLastError()
	}
	ql := createQueueList(cqlist)
	C.drmaa2_list_free(&cqlist)
	return ql, nil
}

// GetAllMachines returns a list of all machines configured in cluster if the argument
// is nil. Otherwise a list of available machines which matches the
// given names is returned.
func (ms *MonitoringSession) GetAllMachines(names []string) (machines []Machine, err error) {
	var arg C.drmaa2_string_list
	if names == nil {
		arg = nil
	} else {
		arg = C.drmaa2_string_list(convertGoListToC(names))
	}
	milist := (C.drmaa2_list)(C.drmaa2_msession_get_all_machines(ms.ms, arg))
	if milist == nil {
		return nil, makeLastError()
	}
	ml := createMachineList(milist)
	C.drmaa2_list_free(&milist)
	return ml, nil
}

func (ms *MonitoringSession) GetAllReservations() (reservations []Reservation, err error) {
	// TODO implement - optional function  (according to DRMAA2 standard)
	return nil, nil
}

// OpenJobSession opens an existing DRMAA2 job sesssion. In Univa Grid Engine
// this job session is persistently stored in the Grid Engine master process.
// The sessionName needs to be != "".
func (sm *SessionManager) OpenJobSession(sessionName string) (*JobSession, error) {
	// convert parameters
	name := C.CString(sessionName)
	defer C.free(unsafe.Pointer(name))
	// DRMAA2 C API call
	var js JobSession
	js.js = C.drmaa2_open_jsession(name)
	// convert error back to Go
	if js.js == nil {
		// an error happended - create an error
		return nil, makeLastError()
	}
	// job session needs to be freed from caller
	return &js, nil
}

// OpenReservationSession opens an existing ReservationSession by name.
func (sm *SessionManager) OpenReservationSession(name string) (rs ReservationSession, err error) {
	return rs, nil
}

/* Destroys an existing session (job session or reservation sesssion). */
func (sm *SessionManager) destroySession(sessionName string, jobSession bool) error {
	// convert parameters
	name := C.CString(sessionName)
	defer C.free(unsafe.Pointer(name))

	var cerror C.drmaa2_error
	// DRMAA2 C API call
	if jobSession {
		cerror = C.drmaa2_destroy_jsession(name)
	} else {
		cerror = C.drmaa2_destroy_rsession(name)
	}
	// convert error back to Go
	if cerror != C.DRMAA2_SUCCESS {
		// an error happended - create an error
		return makeLastError()
	}
	// In case of success nil is returned.
	return nil
}

// DestroyJobSession destroys a job session by name.
func (sm *SessionManager) DestroyJobSession(sessionName string) error {
	return sm.destroySession(sessionName, true)
}

// DestroyReservationSession destroys a reservation by name.
func (sm *SessionManager) DestroyReservationSession(sessionName string) error {
	return sm.destroySession(sessionName, false)
}

func (sm *SessionManager) getSessionNames(jobSession bool) ([]string, error) {
	var name_list C.drmaa2_string_list

	if jobSession {
		name_list = C.drmaa2_get_jsession_names()
	} else {
		name_list = C.drmaa2_get_rsession_names()
	}

	if name_list != nil {
		nl := (C.drmaa2_list)(name_list)
		defer C.drmaa2_list_free(&nl)
		return goStringList(name_list), nil
	}
	return nil, makeLastError()
}

/* GetJobSessionNames returns all job sessions accessable to the user. */
func (sm *SessionManager) GetJobSessionNames() ([]string, error) {
	return sm.getSessionNames(true)
}

/* GetReservationSessionNames returns all reservation sessions accessable to the user. */
func (sm *SessionManager) GetReservationSessionNames() ([]string, error) {
	return sm.getSessionNames(false)
}

/* GetDrmsName returns the name of the Distributed Resource Management System. */
func (sm *SessionManager) GetDrmsName() (string, error) {
	name := C.drmaa2_get_drms_name()
	if name != nil {
		defer C.drmaa2_string_free(&name)
		return C.GoString(name), nil
	}
	return "", makeLastError()
}

/* GetDrmsVersion returns the version of the Distributed Resource Management System. */
func (sm *SessionManager) GetDrmsVersion() (*Version, error) {
	cversion := C.drmaa2_get_drms_version()
	if cversion == nil {
		return nil, makeLastError()
	}
	defer C.drmaa2_version_free(&cversion)

	var version Version
	cmaj := cversion.major
	if cmaj != nil {
		version.Major = C.GoString(cmaj)
	}
	cmin := cversion.minor
	if cmin != nil {
		version.Minor = C.GoString(cmin)
	}
	return &version, nil
}

// Checks whether the DRMAA2 implementation supports an optional
// functionality or not.
func (sm *SessionManager) Supports(c Capability) bool {
	capablilty := capMap[c]
	cres := C.drmaa2_supports(capablilty)
	if cres == C.DRMAA2_TRUE {
		return true
	}
	return false
}

// Event functions

type Event int

const (
	NewState = iota
	Migrated
	AttributeChange
)

type Notification struct {
	Evt         Event    `json:"event"`
	JobId       string   `json:"jobId"`
	SessionName string   `json:"sessionName"`
	State       JobState `json:"jobState"`
}

// A Callback is a function which works on the notification
// struct.
type CallbackFunction func(notification Notification)

// This function is called from C whenever an event happens.
// It is used to forward the event to the Go functions.

// export callbackExecution
func callbackExecution(notify C.drmaa2_notification) {
	// Forward the C notification struct to a Go
	// channel which is subscribed by a coroutine
	// (started by RegisterEventNotification). This
	// coroutine calls all registered callback functions.
}

type EventChannel chan Notification

// TODO(dg) RegisterEventNotification needs to install a C callback
// in the C interface.
func (sm *SessionManager) RegisterEventNotification() (EventChannel, error) {
	// TODO store the callback function
	return nil, nil
}

// TODO(dg) UnregisterEventNotification turns off sending events
// through the EventChannel. Internally a NULL callback is registered
// so that no new events are created.
func (sm *SessionManager) UnregisterEventNotification() error {
	return nil
}

// Close closes the connection to an open JobSession.
func (js *JobSession) Close() error {
	if js.js != nil {
		defer C.drmaa2_jsession_free(&js.js)
	}
	cerr := C.drmaa2_close_jsession(js.js)
	if cerr != C.DRMAA2_SUCCESS {
		return makeLastError()
	}
	/* return nil on success (is easier to handle for the caller) */
	return nil
}

// GetContact returns the contact string of the DRM session.
func (js *JobSession) GetContact() (string, error) {
	contact := C.drmaa2_jsession_get_contact(js.js)
	if contact != nil {
		defer C.drmaa2_string_free(&contact)
		return C.GoString(contact), nil
	}
	return "", makeLastError()
}

// TODO method not needed because session name is part of struct
func (js *JobSession) GetSessionName() (string, error) {
	return js.Name, nil
}

// GetJobCategories returns all job categories specified in the
// job session or globally.
func (js *JobSession) GetJobCategories() ([]string, error) {
	clist := (C.drmaa2_list)(C.drmaa2_jsession_get_job_categories(js.js))
	if clist != nil {
		cl := (C.drmaa2_list)(clist)
		defer C.drmaa2_list_free(&cl)
		return goStringList((C.drmaa2_string_list)(clist)), nil
	}
	return nil, makeLastError()
}

// GetJobs returns a list of all jobs currently attached to the
// given JobSession. If a JobInfo argument unequal nil is given
// then this JobInfo element is used for filtering the result.
// For more details about filtering consider the details of the
// DRMAA2 specificaton.
func (js *JobSession) GetJobs(ji *JobInfo) ([]Job, error) {
	// Create the job filter
	var cji C.drmaa2_jinfo
	if ji != nil {
		cji = convertGoJobInfoToC(*ji)
		defer C.drmaa2_jinfo_free(&cji)
	} else {
		cji = nil
	}
	cjlist := (C.drmaa2_list)(C.drmaa2_jsession_get_jobs(js.js, cji))
	if cjlist == nil {
		return nil, makeLastError()
	}
	jlist := (C.drmaa2_j_list)(cjlist)
	jl := convertCJobListToGo(jlist)
	C.drmaa2_list_free(&cjlist)
	return jl, nil
}

// Converts a C jarray into a Go jarray
func convertJarray(cja C.drmaa2_jarray) (ja ArrayJob) {
	// is this reference needed? better use implementation specific
	ja.aj = cja
	ja.id = C.GoString(cja.id)
	ja.sessionName = C.GoString(cja.session_name)
	ja.jobs = convertCJobListToGo(cja.job_list)
	// jt JobTemplate
	if len(ja.jobs) > 0 {
		ja.jobs[0].GetJobTemplate()
	}
	return ja
}

// Returns a reference to an existing ArrayJob based on the given job
// id. In case of an error the error return value is set to != nil.
func (js *JobSession) GetJobArray(id string) (*ArrayJob, error) {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))

	if jarray := C.drmaa2_jsession_get_job_array(js.js, cid); jarray != nil {
		defer C.drmaa2_jarray_free(&jarray)
		ja := convertJarray(jarray)
		return &ja, nil
	}
	return nil, makeLastError()
}

// RunJob submits a job based on the parameters specified in the JobTemplate
// in the cluster. In case of success it returns a pointer to a Job
// element, which can be used for further processing. In case of an
// error the error return value is set.
func (js *JobSession) RunJob(jt JobTemplate) (*Job, error) {
	// create C.drmaa2_jtemplate and fill in values
	cjtemplate := convertGoJtemplateToC(jt)
	defer C.drmaa2_jtemplate_free(&cjtemplate)

	// set extensions into job template
	setExtensionsIntoCObject(unsafe.Pointer(cjtemplate), jt.ExtensionList)

	if cjob := C.drmaa2_jsession_run_job(js.js, cjtemplate); cjob != nil {
		defer C.drmaa2_j_free(&cjob)
		job := convertCJobToGo(cjob)
		return &job, nil
	}
	return nil, makeLastError()
}

// RunBulkJobs submits a JobTemplate to the cluster as an array job (multiple instances
// of the same job, not neccessarly running a the same point in time).
// It requires a JobTemplate filled out at least with a RemoteCommand.
// The begin, end and step parameters specifying how many array job
// instances are submitted and how the instances are numbered (1,10,1
// denotes 10 array job instances numbered from 1 to 10). The maxParallel
// parameter specifies how many of the array job instances should run
// at parallel as maximum (when resources are contrainted then less
// instances could run).
func (js *JobSession) RunBulkJobs(jt JobTemplate, begin int, end int, step int, maxParallel int) (*ArrayJob, error) {
	cjtemplate := convertGoJtemplateToC(jt)
	if cajob := C.drmaa2_jsession_run_bulk_jobs(js.js, cjtemplate, C.longlong(begin),
		C.longlong(end), C.longlong(step), C.longlong(maxParallel)); cajob != nil {
		defer C.drmaa2_jarray_free(&cajob)
		job := convertCArrayJobToGo(cajob)
		return &job, nil
	}
	return nil, makeLastError()
}

// isStarted determines on which event to wait
func (js *JobSession) waitAny(isStarted bool, jobs []Job, timeout int64) (*Job, error) {
	jl := C.drmaa2_j_list(convertGoListToC(jobs))
	cl := (C.drmaa2_list)(jl)
	defer C.drmaa2_list_free(&cl)

	var cjob C.drmaa2_j
	if isStarted {
		cjob = C.drmaa2_jsession_wait_any_started(js.js, jl, C.time_t(timeout))
	} else {
		cjob = C.drmaa2_jsession_wait_any_terminated(js.js, jl, C.time_t(timeout))
	}
	if cjob != nil {
		job := convertCJobToGo(cjob)
		return &job, nil
	}
	return nil, makeLastError()
}

// WaitAnyStarted waits until any of the given jobs is started (usually in running state).
// The timeout determines after how many seconds the method should abort,
// even when none of the given jobs was started. Special timeout values are
// InfiniteTime and ZeroTime.
func (js *JobSession) WaitAnyStarted(jobs []Job, timeout int64) (*Job, error) {
	return js.waitAny(true, jobs, timeout)
}

// WaitAnyTerminated waits until any of the given jobs is finished. The timeout determines after
// how many seconds the method should abort, even when none of the given jobs
// is finished. Sepecial timeout values are InfiniteTime and ZeroTime.
func (js *JobSession) WaitAnyTerminated(jobs []Job, timeout int64) (*Job, error) {
	return js.waitAny(false, jobs, timeout)
}

// ArrayJob methods.

// GetId returns the job identifier of the ArrayJob.
func (aj *ArrayJob) GetId() string {
	return aj.id
}

// GetJobs returns a list of individual jobs the ArrayJob
// consists of.
func (aj *ArrayJob) GetJobs() []Job {
	return aj.jobs
}

// GetSessionName returns the name of the job session the array job
// belongs to.
func (aj *ArrayJob) GetSessionName() string {
	return aj.sessionName
}

// GetJobTemplate returns a pointer to a job template of the ArrayJob.
func (aj *ArrayJob) GetJobTemplate() *JobTemplate {
	return &aj.jt
}

// Suspend stops all running tasks of an ArrayJob without
// destroying them (usually a SIGSTP signal is send to the
// processes the job consists of).
func (aj *ArrayJob) Suspend() error {
	cjob := convertGoArrayJobToC(*aj)
	defer C.drmaa2_jarray_free(&cjob)
	if C.drmaa2_jarray_suspend(cjob) != C.DRMAA2_SUCCESS {
		return makeLastError()
	}
	return nil
}

// Resume continues all suspended tasks of an ArrayJob.
func (aj *ArrayJob) Resume() error {
	cjob := convertGoArrayJobToC(*aj)
	defer C.drmaa2_jarray_free(&cjob)
	if C.drmaa2_jarray_resume(cjob) != C.DRMAA2_SUCCESS {
		return makeLastError()
	}
	return nil
}

// Hold puts all tasks of an ArrayJob into hold state.
func (aj *ArrayJob) Hold() error {
	cjob := convertGoArrayJobToC(*aj)
	defer C.drmaa2_jarray_free(&cjob)
	if C.drmaa2_jarray_hold(cjob) != C.DRMAA2_SUCCESS {
		return makeLastError()
	}
	return nil
}

// Release releases all tasks of an ArrayJob from hold, if they
// are on hold.
func (aj *ArrayJob) Release() error {
	cjob := convertGoArrayJobToC(*aj)
	defer C.drmaa2_jarray_free(&cjob)
	if C.drmaa2_jarray_release(cjob) != C.DRMAA2_SUCCESS {
		return makeLastError()
	}
	return nil
}

// Terminate terminates (usually sends a SIGKILL) all tasks of an
// ArrayJob.
func (aj *ArrayJob) Terminate() error {
	cjob := convertGoArrayJobToC(*aj)
	defer C.drmaa2_jarray_free(&cjob)
	if C.drmaa2_jarray_terminate(cjob) != C.DRMAA2_SUCCESS {
		return makeLastError()
	}
	return nil
}

// Close closes an open ReservationSession.
func (rs *ReservationSession) Close() error {
	if rs.rs != nil {
		defer C.drmaa2_rsession_free(&rs.rs)
	}
	cerr := C.drmaa2_close_rsession(rs.rs)
	if cerr != C.DRMAA2_SUCCESS {
		return makeLastError()
	}
	return nil
}

// Returns the contact string of the reservation session.
func (rs *ReservationSession) GetContact() (string, error) {
	contact := C.drmaa2_rsession_get_contact(rs.rs)
	if contact != nil {
		defer C.drmaa2_string_free(&contact)
		return C.GoString(contact), nil
	}
	return "", makeLastError()
}

// TODO(dg) Returns the name of the reservation session
func (rs *ReservationSession) GetSessionName() (string, error) {
	return rs.Name, nil
}

// TODO(dg) Returns a reservation object based on the AR id
func (rs *ReservationSession) GetReservation(rid string) (*Reservation, error) {
	return nil, nil
}

// TODO(dg) Allocates an advance reservation based on the reservation template
func (rs *ReservationSession) RequestReservation(rtemplate ReservationTemplate) (*Reservation, error) {
	return nil, nil
}

// TODO(dg) Returns a list of available advance reservations
func (rs *ReservationSession) GetReservations() ([]Reservation, error) {
	// TODO implement
	return nil, nil
}

// TODO(dg) Returns the advance reservation id
func (r *Reservation) GetId() (string, error) {
	// TODO implement
	return "", nil
}

// TODO(dg) Returns the name of the reservation
func (r *Reservation) GetSessionName() (string, error) {
	// TODO implement
	return "", nil
}

// TODO(dg) Returns the reservation template of the reservation
func (r *Reservation) GetTemplate() (*ReservationTemplate, error) {
	// TODO implement
	return nil, nil
}

// TODO(dg) Returns the reservation info object of the reservation
func (r *Reservation) GetInfo() (*ReservationInfo, error) {
	// TODO implement
	return nil, nil
}

// TODO(dg) Cancels an advance reservation
func (r *Reservation) Terminate() error {
	// TODO implement
	return nil
}

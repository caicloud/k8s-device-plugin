package v1alpha2

// Framework type of application
type FrameworkType string

// Framework defines include machine learning framework and programming language
const (
	// Algorithm framework
	Tensorflow FrameworkType = "tensorflow"
	TFserving  FrameworkType = "tfserving"
	Pytorch    FrameworkType = "pytorch"
	SKLearn    FrameworkType = "sklearn"
	Chainer    FrameworkType = "chainer"
	Caffe2     FrameworkType = "caffe2"
	Caffe      FrameworkType = "caffe"
	MXNet      FrameworkType = "mxnet"
	Spark      FrameworkType = "spark"
	Keras      FrameworkType = "keras"
	Onnx       FrameworkType = "onnx"
	Horovod    FrameworkType = "horovod"

	// Programming language
	JavaScript FrameworkType = "javascript"
	Python     FrameworkType = "python"
	Golang     FrameworkType = "golang"
	Scala      FrameworkType = "scala"
	Java       FrameworkType = "java"
	Cpp        FrameworkType = "cpp"
	C          FrameworkType = "c"
	R          FrameworkType = "r"

	// Programming development tools
	Jupyter            FrameworkType = "jupyter"
	Zeppelin           FrameworkType = "zeppelin"
	JupyterLab         FrameworkType = "jupyterlab"
	TensorBoard        FrameworkType = "tensorboard"
	SparkHistoryServer FrameworkType = "sparkhistoryserver"
)

// Framework group type of framework
type FrameworkGroupType string

// Framework group defines
const (
	FeatureEnginnering = "FeatureEngineering"
	DevelopmentTools   = "DevelopmentTools"
	ParameterTuning    = "ParameterTuning"
	ModelEvaluation    = "ModelEvaluation"
	ModelTraining      = "ModelTraining"
	ModelServing       = "ModelServing"
	DataCleaning       = "DataCleaning"
)

type DataType string

const (
	// Model set type
	Model   DataType = "model"
	Dataset DataType = "dataset"
)

// ReplicaType represents the type of the replica. Each operator needs to define its
// own set of ReplicaTypes.
type MLNeuronReplicaType string

type SparkReplicaType MLNeuronReplicaType

const (
	// SparkReplicaTypeDriver is the type of Driver in Spark
	SparkReplicaTypeDriver SparkReplicaType = "Driver"

	// SparkReplicaTypeExecutor is the type for Executor in Spark.
	SparkReplicaTypeExecutor SparkReplicaType = "Executor"
)

// TFReplicaType is the type for TFReplica.
type TFReplicaType MLNeuronReplicaType

const (
	// TFReplicaTypePS is the type for parameter servers of distributed TensorFlow.
	TFReplicaTypePS TFReplicaType = "PS"

	// TFReplicaTypeWorker is the type for workers of distributed TensorFlow.
	// This is also used for non-distributed TensorFlow.
	TFReplicaTypeWorker TFReplicaType = "Worker"

	// TFReplicaTypeChief is the type for chief worker of distributed TensorFlow.
	// If there is "chief" replica type, it's the "chief worker".
	// Else, worker:0 is the chief worker.
	TFReplicaTypeChief TFReplicaType = "Chief"

	// TFReplicaTypeMaster is the type for master worker of distributed TensorFlow.
	// This is similar to chief, and kept just for backwards compatibility.
	TFReplicaTypeMaster TFReplicaType = "Master"

	// TFReplicaTypeEval is the type for evaluation replica in TensorFlow.
	TFReplicaTypeEval TFReplicaType = "Evaluator"
)

// PyTorchReplicaType is the type for PyTorchReplica.
type PyTorchReplicaType MLNeuronReplicaType

const (
	// PyTorchReplicaTypeMaster is the type of Master of distributed PyTorch
	PyTorchReplicaTypeMaster PyTorchReplicaType = "Master"

	// PyTorchReplicaTypeWorker is the type for workers of distributed PyTorch.
	PyTorchReplicaTypeWorker PyTorchReplicaType = "Worker"
)

// MXReplicaType is the type for MXReplica.
type MXReplicaType MLNeuronReplicaType

const (
	// MXReplicaTypeScheduler is the type for scheduler replica in MXNet.
	MXReplicaTypeScheduler MXReplicaType = "Scheduler"

	// MXReplicaTypeServer is the type for parameter servers of distributed MXNet.
	MXReplicaTypeServer MXReplicaType = "Server"

	// MXReplicaTypeWorker is the type for workers of distributed MXNet.
	// This is also used for non-distributed MXNet.
	MXReplicaTypeWorker MXReplicaType = "Worker"

	// MXReplicaTypeTunerTracker
	// This the auto-tuning tracker e.g. autotvm tracker, it will dispatch tuning task to TunerServer
	MXReplicaTypeTunerTracker MXReplicaType = "TunerTracker"

	// MXReplicaTypeTunerServer
	MXReplicaTypeTunerServer MXReplicaType = "TunerServer"

	// MXReplicaTuner is the type for auto-tuning of distributed MXNet.
	// This is also used for non-distributed MXNet.
	MXReplicaTypeTuner MXReplicaType = "Tuner"
)

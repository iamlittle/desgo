# Desgo

Desgo is a multithreaded [discrete event simulation](https://en.wikipedia.org/wiki/Discrete_event_simulation) cashier example written in Go. Some rudimentary output analysis done in Jupyter notebooks accompanies. Currently, it is a single cashier / customer simulation engine but I would like to expand this to include a warehouse inventory / fill rate engine.

### Building the cashier executable

Desgo requires that you have [Go installed](https://golang.org/doc/install). 
 
```
cd github.com/iamlittle/desgo/cashier
go build
```

### Running an example simulation

```
cd github.com/iamlittle/desgo/
cashier/cashier --input ./exploration/cashier/01_initial_look/input.01.yaml
```

### Configuring a simulation

Each simulation is configured using a yaml file. [Here is an example](./exploration/cashier/01_initial_look/input.01.yaml). Here you can configure the number of iterations run, number of customers, cashiers, distribution parameters, etc.

### Jupyter Notebooks

See some rudimentary analysis in the exploration steps [Initial Look](./exploration/cashier/01_initial_look/01_initial_look.ipynb) and [Multiple Runs](./exploration/02_multiple_runs/02_multiple_runs.ipynb)


 
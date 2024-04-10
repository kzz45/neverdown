/*
Copyright 2017 The Kubernetes Authors.

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

package internalversion

import (
	"github.com/kzz45/neverdown/pkg/kubernetes/printers"
)

// AddHandlers adds print handlers for default Kubernetes types dealing with internal versions.
func AddHandlers(h printers.PrintHandler) {

}

// func printCoffeeList(coffeeList *v1.CoffeeList, options printers.GenerateOptions) ([]metav1.TableRow, error) {
// 	rows := make([]metav1.TableRow, 0, len(coffeeList.Items))
// 	for i := range coffeeList.Items {
// 		r, err := printCoffee(&coffeeList.Items[i], options)
// 		if err != nil {
// 			return nil, err
// 		}
// 		rows = append(rows, r...)
// 	}
// 	return rows, nil
// }

// func printCoffee(coffee *v1.Coffee, options printers.GenerateOptions) ([]metav1.TableRow, error) {
// 	row := metav1.TableRow{
// 		Object: runtime.RawExtension{Object: coffee},
// 	}
// 	beans := coffee.Spec.Beans
// 	row.Cells = append(row.Cells, coffee.Name, beans)
// 	return []metav1.TableRow{row}, nil
// }

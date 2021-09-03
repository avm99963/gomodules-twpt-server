// We just import the components needed.
import {button, checkbox, dataTable, dialog, drawer, fab, iconButton, list, materialIcon, menu, select, textfield, topAppBar} from 'vue-material-adapter';

export default {
  install(vm) {
    vm.use(button);
    vm.use(checkbox);
    vm.use(dataTable);
    vm.use(dialog);
    vm.use(drawer);
    vm.use(fab);
    vm.use(iconButton);
    vm.use(list);
    vm.use(materialIcon);
    vm.use(menu);
    vm.use(select);
    vm.use(textfield);
    vm.use(topAppBar);
  },
};

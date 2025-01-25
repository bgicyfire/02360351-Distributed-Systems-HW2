import {Routes} from '@angular/router';
import {ServersListComponent} from './components/servers-list/servers-list.component';
import {ScootersListComponent} from './components/scooters-list/scooters-list.component';

export const routes: Routes = [
  {
    pathMatch: 'prefix',
    path: '',
    redirectTo: 'scooters'
  },
  {
    path: 'scooters',
    component: ScootersListComponent
  },
  {
    path: 'servers',
    component: ServersListComponent
  }
];

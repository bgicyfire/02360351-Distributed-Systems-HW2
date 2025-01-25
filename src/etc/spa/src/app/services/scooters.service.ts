import {Injectable} from '@angular/core';
import {map, Observable} from 'rxjs';
import {HttpClient} from '@angular/common/http';
import {Scooter} from '../dtos/scooter';

@Injectable({
  providedIn: 'root'
})
export class ScootersService {
  private BASE_URL = 'http://localhost:50053';

  constructor(private http: HttpClient) {
  }

  public getServersList(): Observable<string[]> {
    return this.http.get<{ servers: string[] }>(this.BASE_URL + '/servers')
      .pipe(map(a => a.servers));
  }

  public getScootersList(): Observable<Scooter[]> {
    return this.http.get<Scooter[]>(this.BASE_URL + '/scooters').pipe(
      map(scooters => scooters.sort((a, b) => a.id.localeCompare(b.id)))
    );
  }

  public createScooter(id: string): Observable<Scooter> {
    let data: Scooter = {
      'id': id,
      'is_available': true,
      'total_distance': 0,
      'current_reservation_id': ''
    };
    return this.http.put<Scooter>(this.BASE_URL + '/scooters/' + id, data);
  }

  public reserveScooter(scooterId: string): Observable<string> {
    return this.http.post<string>(this.BASE_URL + '/scooters/' + scooterId + '/reservations', {});
  }

  public releaseScooter(scooterId: string, reservationId: string, rideDistance: number): Observable<any> {
    // TODO: decide if also want to require the reservation id
    let data = {
      scooterId: scooterId,
      reservation_id: reservationId,
      ride_distance: rideDistance
    }
    return this.http.post<any>(this.BASE_URL + '/scooters/' + scooterId + '/releases', data);
  }
}

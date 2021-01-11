import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ClusterHealthCheckComponent } from './cluster-health-check.component';

describe('ClusterHealthCheckComponent', () => {
  let component: ClusterHealthCheckComponent;
  let fixture: ComponentFixture<ClusterHealthCheckComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ClusterHealthCheckComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ClusterHealthCheckComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

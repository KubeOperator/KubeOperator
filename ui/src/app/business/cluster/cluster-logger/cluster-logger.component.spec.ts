import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ClusterLoggerComponent } from './cluster-logger.component';

describe('ClusterLoggerComponent', () => {
  let component: ClusterLoggerComponent;
  let fixture: ComponentFixture<ClusterLoggerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ClusterLoggerComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ClusterLoggerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

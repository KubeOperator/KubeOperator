import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ConfigMapListComponent } from './config-map-list.component';

describe('ConfigMapListComponent', () => {
  let component: ConfigMapListComponent;
  let fixture: ComponentFixture<ConfigMapListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ConfigMapListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ConfigMapListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
